package cuckoo

import (
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestAeternity(t *testing.T) {
	tests := []struct {
		header []byte
		sols   []uint64
	}{
		{
			header: testutil.MustDecodeHex("69435549583467653534384d6e4e6c4d714b6371365444674a6631446b783548594f4d67705041783249343d71414541414837384141413d000000000000000000000000000000000000000000000000"),
			sols: []uint64{
				0x003b5d47, 0x00a70508, 0x00d0aa4a, 0x0238a16a, 0x038653bf, 0x03e91d96, 0x03f4baa8, 0x062ef17e,
				0x065d7b41, 0x066fbb1e, 0x079af861, 0x08bd2cf2, 0x0956b89d, 0x0b56fb7f, 0x0c098553, 0x0c6d2c27,
				0x0d8c0fd9, 0x0ddcbb1d, 0x0e3eccde, 0x0e464bef, 0x0fb09bef, 0x1267ebb1, 0x129ef8e6, 0x138432b5,
				0x144d428b, 0x1484e6b6, 0x14efcfba, 0x158d5352, 0x159f3551, 0x15a07563, 0x160a3efd, 0x17c9b61e,
				0x184499bc, 0x1844f434, 0x1919053a, 0x197a9095, 0x1aa04947, 0x1bc3f6e5, 0x1d8b4029, 0x1e6a1fe0,
				0x1e7e4380, 0x1f5a2a50,
			},
		},
		{
			header: testutil.MustDecodeHex("5957616c3675576330543433356d6775473877482f4d4378354e45626f6a4837586d38547176614e3862553d566b35374733496f4141413d000000000000000000000000000000000000000000000000"),
			sols: []uint64{
				0x007791c9, 0x0197cfec, 0x025166a5, 0x025e30b8, 0x03d6aedc, 0x04311736, 0x0516b291, 0x05179816,
				0x056ac771, 0x068b92d8, 0x06aee4f8, 0x0817c500, 0x08c1a751, 0x08f1999c, 0x08f46654, 0x092b8770,
				0x096034f8, 0x09cfcb6f, 0x0d4f1cdb, 0x0d5b8793, 0x0d5c09a8, 0x0da0e97b, 0x0db38171, 0x0dc09c91,
				0x0e34aeaf, 0x0f0251df, 0x10486aa7, 0x14477924, 0x1544daf9, 0x1623e36e, 0x16b55699, 0x16fa6905,
				0x177125c6, 0x17f88a4d, 0x19243b3c, 0x19ebdaa7, 0x1b5383f1, 0x1c614cbf, 0x1de144f9, 0x1df596a7,
				0x1f629453, 0x1f6b9406,
			},
		},
		{
			header: testutil.MustDecodeHex("6a4f3142484c734a775a6571367430457650794d76394a44583274387162574859586e496171494d3245773d582b54766c6e496f4141413d000000000000000000000000000000000000000000000000"),
			sols: []uint64{
				0x008a9cad, 0x008ec69d, 0x00dbf03e, 0x020337ed, 0x0211161e, 0x02ee6658, 0x045d4488, 0x0534f601,
				0x0542e80c, 0x062d2002, 0x078ac9b5, 0x07b7112a, 0x07cb2202, 0x07fbc7e7, 0x088827a7, 0x093e3139,
				0x0a7af764, 0x0a92c274, 0x0ad8dcb9, 0x0bf8f10f, 0x0c0f0de0, 0x0cc6670d, 0x0dad83b3, 0x104f3b5e,
				0x11817d4c, 0x123a7b96, 0x12b9fbdd, 0x1337ce2b, 0x13432b38, 0x15b1b455, 0x16784bf4, 0x1693250e,
				0x169d804f, 0x16e811de, 0x195d114e, 0x1b77f0cb, 0x1c55fca6, 0x1cab0177, 0x1cb0a8c3, 0x1d38ea27,
				0x1d9aa6dc, 0x1f713d1f,
			},
		},
	}

	for i, tt := range tests {
		valid, err := NewAeternity().Verify(tt.header, tt.sols)
		if err != nil {
			t.Errorf("failed on %d: %v", i, err)
		} else if !valid {
			t.Errorf("failed on %d: invalid solution", i)
		}
	}
}

func TestCortex(t *testing.T) {
	tests := []struct {
		header []byte
		sols   []uint64
	}{
		{
			header: testutil.MustDecodeHex("6281a031a95a7669e42cf56d46b5d921b067ace29c46c89fa2698f3b895d6fcb21208e4e00000165"),
			sols: []uint64{
				0x017ca085, 0x0181ca71, 0x096b8b98, 0x09d3a607, 0x0b6bb4c8, 0x0c9bbecb, 0x10d1c645, 0x13ba80dc,
				0x13cb4dc9, 0x15ebc37d, 0x164de862, 0x16a7906a, 0x18c28113, 0x199e50ca, 0x1ba70932, 0x1bc435b1,
				0x1caad714, 0x1d94ccd4, 0x1da4b49d, 0x1eff189e, 0x2030c2cf, 0x2084a6c3, 0x2111e51e, 0x241ff2d0,
				0x26bb0111, 0x275fd4a1, 0x27654850, 0x291041de, 0x2a4c1e5b, 0x2a8e54e1, 0x2ba12d29, 0x2d16cbc0,
				0x2e9e0df8, 0x3209259d, 0x32751e22, 0x33107850, 0x332b35f9, 0x33a134d4, 0x354fc224, 0x384052fb,
				0x38cdb22e, 0x3e665fed,
			},
		},
		{
			header: testutil.MustDecodeHex("a855bfb05489f16720c241b8f89db047fd63b4e7ffc8942c37edd81c6a4ee2e321208e4e0000022a"),
			sols: []uint64{
				0x005f7571, 0x008178c5, 0x00b26db9, 0x01045303, 0x013a9f66, 0x021e09af, 0x0223b428, 0x024fe83b,
				0x03f17913, 0x03f232d0, 0x050e8553, 0x054dcf61, 0x05db04d4, 0x0696f89c, 0x0834e12e, 0x0857a4d3,
				0x0a2fd97b, 0x0aa26d69, 0x0ce78e0b, 0x0d468cfc, 0x0d5413af, 0x0d6fa7d1, 0x0d91d094, 0x0f173576,
				0x1450b48b, 0x182aa6fb, 0x1833ba9a, 0x1a06e9c8, 0x1a2f0366, 0x1c7eefcf, 0x1eb26af8, 0x20021c3a,
				0x216dba37, 0x2597bce1, 0x25a187a8, 0x26763784, 0x2a4f0d5e, 0x2d48073f, 0x2f3581b5, 0x3518197e,
				0x37f4fb08, 0x3d05d79f,
			},
		},
		{
			header: testutil.MustDecodeHex("a855bfb05489f16720c241b8f89db047fd63b4e7ffc8942c37edd81c6a4ee2e321208e4e00000328"),
			sols: []uint64{
				0x001b27e9, 0x003af2c7, 0x01eed987, 0x039a8920, 0x050de2d0, 0x07f4dde9, 0x07fd54af, 0x0c4e256a,
				0x0dbfe2b4, 0x0ef70428, 0x0f103eb0, 0x0f3d3c08, 0x103f426f, 0x14f152bd, 0x15f55244, 0x16531636,
				0x16bce128, 0x16ceb3d8, 0x1774ebef, 0x1b25e464, 0x1d837e14, 0x24fcf2ec, 0x264e7881, 0x284f5e09,
				0x2994e92d, 0x2cab769d, 0x2d85985f, 0x2dbc88f1, 0x2f3279b1, 0x2f5984a2, 0x30480ed6, 0x30f44f85,
				0x32452faf, 0x33b1a793, 0x33ea0ab1, 0x3406943e, 0x34587460, 0x34c301cd, 0x36bd7591, 0x386e9b94,
				0x39e03a63, 0x3a7c54c3,
			},
		},
	}

	for i, tt := range tests {
		valid, err := NewCortex().Verify(tt.header, tt.sols)
		if err != nil {
			t.Errorf("failed on %d: %v", i, err)
		} else if !valid {
			t.Errorf("failed on %d: invalid solution", i)
		}
	}
}
