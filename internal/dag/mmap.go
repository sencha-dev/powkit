// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package dag

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

var (
	dumpMagic = []uint32{0xbaddcafe, 0xfee1dead}
)

type dataFile struct {
	dump *os.File
	mmap mmap.MMap
	data []uint32
}

// memoryMap tries to memory map a file of uint32s for read only access.
func memoryMap(path string, lock bool) (dataFile, error) {
	var df dataFile

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return df, err
	}

	mem, buffer, err := memoryMapFile(file, false)
	if err != nil {
		file.Close()
		return df, err
	}

	for i, magic := range dumpMagic {
		if buffer[i] != magic {
			mem.Unmap()
			file.Close()
			return df, errors.New("invalid dump magic")
		}
	}

	if lock {
		if err := mem.Lock(); err != nil {
			mem.Unmap()
			file.Close()
			return df, errors.New("invalid dump magic")
		}
	}

	df = dataFile{
		dump: file,
		mmap: mem,
		data: buffer[len(dumpMagic):],
	}

	return df, nil
}

// memoryMapFile tries to memory map an already opened file descriptor.
func memoryMapFile(file *os.File, write bool) (mmap.MMap, []uint32, error) {
	// Try to memory map the file
	flag := mmap.RDONLY
	if write {
		flag = mmap.RDWR
	}

	mem, err := mmap.Map(file, flag, 0)
	if err != nil {
		return nil, nil, err
	}

	// The file is now memory-mapped. Create a []uint32 view of the file.
	var view []uint32
	header := (*reflect.SliceHeader)(unsafe.Pointer(&view))
	header.Data = (*reflect.SliceHeader)(unsafe.Pointer(&mem)).Data
	header.Cap = len(mem) / 4
	header.Len = header.Cap
	return mem, view, nil
}

// ensureSize expands the file to the given size. This is to prevent runtime
// errors later on, if the underlying file expands beyond the disk capacity,
// even though it ostensibly is already expanded, but due to being sparse
// does not actually occupy the full declared size on disk.
func ensureSize(f *os.File, size int64) error {
	// On systems which do not support fallocate, we merely truncate it.
	// More robust alternatives  would be to
	// - Use posix_fallocate, or
	// - explicitly fill the file with zeroes.
	return f.Truncate(size)
}

// memoryMapAndGenerate tries to memory map a temporary file of uint32s for write
// access, fill it with the data from a generator and then move it into the final
// path requested.
func memoryMapAndGenerate(path string, size uint64, lock bool, generator func(buffer []uint32)) (dataFile, error) {
	var df dataFile

	// Ensure the data folder exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return df, err
	}

	// Create a huge temporary empty file to fill with data
	temp := path + "." + strconv.Itoa(rand.Int())
	dump, err := os.Create(temp)
	if err != nil {
		return df, err
	}

	if err = ensureSize(dump, int64(len(dumpMagic))*4+int64(size)); err != nil {
		dump.Close()
		os.Remove(temp)
		return df, err
	}

	// Memory map the file for writing and fill it with the generator
	mem, buffer, err := memoryMapFile(dump, true)
	if err != nil {
		dump.Close()
		os.Remove(temp)
		return df, err
	}

	copy(buffer, dumpMagic)
	data := buffer[len(dumpMagic):]
	generator(data)

	if err := mem.Unmap(); err != nil {
		return df, err
	}

	if err := dump.Close(); err != nil {
		return df, err
	}

	if err := os.Rename(temp, path); err != nil {
		return df, err
	}

	return memoryMap(path, lock)
}
