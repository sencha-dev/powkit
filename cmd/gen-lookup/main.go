package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/sencha-dev/powkit/internal/dag"
)

func genCacheLookupTable(initBytes, growthBytes uint64, epochs int) []uint64 {
	d := &dag.DAG{
		Config: dag.Config{
			CacheInitBytes:   initBytes,
			CacheGrowthBytes: growthBytes,
		},
	}

	lookupTable := make([]uint64, epochs)
	for i := range lookupTable {
		lookupTable[i] = d.CacheSize(uint64(i))
	}

	return lookupTable
}

func genDatasetLookupTable(mixBytes, initBytes, growthBytes uint64, epochs int) []uint64 {
	d := &dag.DAG{
		Config: dag.Config{
			MixBytes:           mixBytes,
			DatasetInitBytes:   initBytes,
			DatasetGrowthBytes: growthBytes,
		},
	}

	lookupTable := make([]uint64, epochs)
	for i := range lookupTable {
		lookupTable[i] = d.DatasetSize(uint64(i))
	}

	return lookupTable
}

func formatLookupTable(lookupTable []uint64, rowSize int) string {
	rows := len(lookupTable) / rowSize

	var formatted string
	for i := 0; i < rows; i++ {
		var row string
		for j := 0; j < rowSize; j++ {
			index := (i * rowSize) + j
			if index >= len(lookupTable) {
				continue
			}

			row += fmt.Sprintf("%d,", lookupTable[index])
			if j < rowSize-1 {
				row += " "
			}
		}

		formatted += fmt.Sprintf("\t%s\n", row)
	}

	val := fmt.Sprintf("[]uint64{\n%s}", formatted)

	return val
}

func saveToFile(data []byte, name string) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("locating working directory: %s", err)
	}

	filename := fmt.Sprintf("%s_lookup.go", name)
	path := filepath.Join(dir, filename)

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

var lookupTemplate = `package {{.Package}}

var cacheSizes = {{.CacheLookup}}

var datasetSizes = {{.DatasetLookup}}
`

func main() {
	var packageName string
	var epochs int
	var mixBytes, cacheInitBytes, cacheGrowthBytes, datasetInitBytes, datasetGrowthBytes uint64

	flag.StringVar(&packageName, "package", "", "The package name where the lookup table will be generated")
	flag.IntVar(&epochs, "epochs", 2048, "The number of epochs to generate for")
	flag.Uint64Var(&mixBytes, "mixBytes", 128, "The number of bytes in a mix hash")
	flag.Uint64Var(&cacheInitBytes, "cacheInit", 0, "The cache initialization size in bytes")
	flag.Uint64Var(&cacheGrowthBytes, "cacheGrowth", 0, "The cache growth factor in bytes")
	flag.Uint64Var(&datasetInitBytes, "datasetInit", 0, "The dataset initialization size in bytes")
	flag.Uint64Var(&datasetGrowthBytes, "datasetGrowth", 0, "The dataset growth factor in bytes")
	flag.Parse()

	if len(packageName) == 0 {
		log.Fatalf("package name is required")
	}

	cacheLookupTable := genCacheLookupTable(cacheInitBytes, cacheGrowthBytes, epochs)
	datasetLookupTable := genDatasetLookupTable(mixBytes, datasetInitBytes, datasetGrowthBytes, epochs)

	data := map[string]interface{}{
		"Package":       packageName,
		"CacheLookup":   formatLookupTable(cacheLookupTable, 5),
		"DatasetLookup": formatLookupTable(datasetLookupTable, 5),
	}

	buf := new(bytes.Buffer)
	t := template.Must(template.New("lookup").Parse(lookupTemplate))
	t.Execute(buf, data)

	saveToFile(buf.Bytes(), packageName)
}
