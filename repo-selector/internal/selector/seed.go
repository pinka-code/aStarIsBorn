package selector

import (
	"crypto/sha1"
	"encoding/binary"
	"time"
)

func SeedFromDate(date time.Time) int64 {
	h := sha1.Sum([]byte(date.Format("2006-01-02")))
	return int64(binary.BigEndian.Uint64(h[:8]))
}