package sqlxmemo

import (
	"encoding/hex"
	gocrc32 "hash/crc32"
)

var crc32poly = gocrc32.MakeTable(gocrc32.Koopman)

// FromString converts string into CRC32 hash string.
func FromString(str string) string {
	hash := gocrc32.New(crc32poly)
	_, _ = hash.Write([]byte(str))
	// Encode the hash to a string
	return hex.EncodeToString(hash.Sum(nil))
}
