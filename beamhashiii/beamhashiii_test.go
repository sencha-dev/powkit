package beamhashiii

import (
	"testing"
)

func TestVerifyBeam(t *testing.T) {
	tests := []struct {
		header []byte
		soln   []byte
	}{
		{
			header: []byte{
				0xfc, 0x40, 0x99, 0x6a, 0x51, 0x8c, 0x22, 0x13,
				0x84, 0xc9, 0xf2, 0x54, 0x2c, 0xa8, 0x11, 0xcd,
				0x66, 0xc4, 0xcc, 0xdd, 0xb0, 0x01, 0xef, 0x40,
				0xb9, 0xf9, 0xba, 0x05, 0x9c, 0x20, 0x35, 0x2e,
				0xb3, 0x2c, 0x7d, 0x4f, 0x07, 0xa3, 0x00, 0x1c,
			},
			soln: []byte{
				0x0f, 0xc8, 0x1c, 0x68, 0x4b, 0xe2, 0x29, 0xc3,
				0x6b, 0x84, 0x4e, 0xf8, 0x29, 0x9a, 0x97, 0x44,
				0xdb, 0xb8, 0x72, 0x72, 0x76, 0xbf, 0xf8, 0xcb,
				0xd6, 0x10, 0xfa, 0x74, 0x14, 0xfb, 0x6c, 0xfd,
				0x67, 0xb9, 0x25, 0x86, 0xf8, 0x4f, 0x8b, 0xff,
				0xae, 0xeb, 0x99, 0x26, 0x69, 0x94, 0xd7, 0x9d,
				0xa3, 0xfb, 0x02, 0x6a, 0x24, 0x12, 0x8b, 0x84,
				0x90, 0x1f, 0x24, 0x4b, 0x08, 0xee, 0x6b, 0x6b,
				0x95, 0x43, 0x72, 0xfc, 0xb0, 0xa7, 0xd3, 0x33,
				0x18, 0xda, 0x6b, 0xf1, 0x85, 0x4a, 0xe4, 0x8f,
				0x94, 0xfe, 0x8a, 0xf2, 0xd3, 0x14, 0x7b, 0xdc,
				0x73, 0x02, 0xcc, 0x12, 0xda, 0xa1, 0xa3, 0x06,
				0x51, 0x11, 0x22, 0xa7, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			header: []byte{
				0x0a, 0x24, 0xff, 0x76, 0xe1, 0xf5, 0x3d, 0xdf,
				0x60, 0x5c, 0x66, 0xa2, 0x1e, 0x9e, 0x67, 0xa6,
				0x8f, 0xd3, 0x7e, 0xcd, 0xec, 0x17, 0xda, 0xd8,
				0x50, 0x02, 0x97, 0x62, 0x2f, 0x9d, 0x9a, 0xd4,
				0x8f, 0x92, 0xd1, 0x19, 0xdf, 0x9a, 0x53, 0x46,
			},
			soln: []byte{
				0x55, 0x98, 0x39, 0xf6, 0x84, 0xc3, 0xbd, 0x54,
				0xbe, 0x26, 0x5d, 0xc3, 0xcf, 0x50, 0xc6, 0x6c,
				0x08, 0x4a, 0x61, 0xe9, 0xaa, 0xca, 0x8c, 0x1e,
				0xca, 0x5e, 0xb8, 0xc6, 0xbe, 0x10, 0xd9, 0x1e,
				0xc6, 0x00, 0x2d, 0x58, 0x01, 0x8c, 0x19, 0x98,
				0x4c, 0x47, 0x66, 0x9a, 0x28, 0x87, 0x4c, 0xb6,
				0x37, 0xec, 0x59, 0x95, 0x64, 0x24, 0xd2, 0x3b,
				0x1b, 0xba, 0x57, 0xbb, 0x7f, 0xf9, 0x3c, 0xcd,
				0x5f, 0x46, 0x17, 0xd4, 0xab, 0xb5, 0x5f, 0xc0,
				0x1b, 0x15, 0xce, 0xd2, 0x88, 0x66, 0xfc, 0x6c,
				0x45, 0x82, 0xb8, 0x5f, 0x94, 0x22, 0x52, 0x4c,
				0xf5, 0x3e, 0x6f, 0x36, 0xc2, 0x6b, 0xe6, 0x75,
				0xf8, 0x9e, 0xe3, 0xf2, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			header: []byte{
				0xf2, 0x48, 0x52, 0xb0, 0x58, 0x0a, 0xa4, 0x48,
				0xe0, 0xcb, 0xed, 0x5e, 0xe8, 0x4d, 0x23, 0x5c,
				0xe6, 0x22, 0xb8, 0xe4, 0x43, 0x56, 0x1d, 0xbc,
				0xf1, 0x96, 0x14, 0x8d, 0xf2, 0xfc, 0x02, 0x4e,
				0x8f, 0x92, 0xd1, 0x19, 0xf5, 0x97, 0xee, 0xa8,
			},
			soln: []byte{
				0x8d, 0xd2, 0x0d, 0xcc, 0x03, 0xd0, 0x33, 0xa2,
				0x23, 0x0a, 0x87, 0x1f, 0x8e, 0xc1, 0x0a, 0x73,
				0xad, 0x34, 0x75, 0xe5, 0x7e, 0x4f, 0xdc, 0x00,
				0xfe, 0xdb, 0x58, 0x51, 0xa4, 0x5b, 0x79, 0x45,
				0xa8, 0x88, 0x8c, 0xc1, 0x9e, 0xd9, 0xe2, 0xdf,
				0xce, 0x8e, 0xb9, 0xb1, 0xd0, 0x37, 0xc5, 0xdd,
				0x8c, 0xc6, 0x6a, 0xed, 0x27, 0x50, 0x19, 0xf1,
				0xfe, 0x76, 0x5b, 0x89, 0x44, 0xf8, 0x55, 0xcc,
				0x05, 0x0e, 0x70, 0x80, 0xfd, 0xdf, 0x74, 0xd5,
				0x52, 0x68, 0xfb, 0x7d, 0xf2, 0x2d, 0x28, 0xf0,
				0xe5, 0xd6, 0x46, 0x16, 0xa2, 0x91, 0xe3, 0xab,
				0xb8, 0x42, 0xd6, 0xf8, 0xee, 0xbd, 0x4b, 0xfa,
				0xe7, 0x19, 0xab, 0xd8, 0x00, 0x00, 0x00, 0x00,
			},
		},
	}

	for i, tt := range tests {
		valid, err := NewBeam().Verify(tt.header, tt.soln)
		if err != nil {
			t.Errorf("failed on %d: %v", i, err)
		} else if !valid {
			t.Errorf("failed on %d: invalid solution", i)
		}
	}
}
