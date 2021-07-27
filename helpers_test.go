package librsync

import (
	"fmt"
	"strconv"
	"strings"
)

var allTestCases = []string{
	"000-blake2-11-23",
	"000-blake2-512-32",
	"000-md4-256-7",
	"001-blake2-512-32",
	"001-blake2-776-31",
	"001-md4-777-15",
	"002-blake2-512-32",
	"002-blake2-431-19",
	"002-md4-128-16",
	"003-blake2-512-32",
	"003-blake2-1024-13",
	"003-md4-1024-13",
	"004-blake2-1024-28",
	"004-blake2-2222-31",
	"004-blake2-512-32",
	"005-blake2-512-32",
	"005-blake2-1000-18",
	"005-md4-999-14",
	"006-blake2-2-32",
	"007-blake2-5-32",
	"007-blake2-4-32",
	"007-blake2-3-32",
	"008-blake2-222-30",
	"008-blake2-512-32",
	"008-md4-111-11",
	"009-blake2-2048-26",
	"009-blake2-512-32",
	"009-md4-2033-15",
	"010-blake2-512-32",
	"010-blake2-7-6",
	"010-md4-4096-8",
}

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
