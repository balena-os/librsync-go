package librsync

import (
	"fmt"
	"strconv"
	"strings"
)

func argsFromTestName(name string) (file string, magic MagicNumber, blockLen, strongLen uint32, err error) {
	segs := strings.Split(name, "-")
	if len(segs) != 4 {
		return "", 0, 0, 0, fmt.Errorf("invalid format for name %q", name)
	}

	file = segs[0]

	switch segs[1] {
	case "blake2":
		magic = BLAKE2_SIG_MAGIC
	case "md4":
		magic = MD4_SIG_MAGIC
	default:
		return "", 0, 0, 0, fmt.Errorf("invalid magic %q", segs[1])
	}

	blockLen64, err := strconv.ParseInt(segs[2], 10, 32)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("invalid block length %q", segs[2])
	}
	blockLen = uint32(blockLen64)

	strongLen64, err := strconv.ParseInt(segs[3], 10, 32)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("invalid strong hash length %q", segs[3])
	}
	strongLen = uint32(strongLen64)

	return
}
