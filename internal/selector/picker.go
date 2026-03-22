package selector

import (
	"math/rand"
)

func PickIndex(seed int64, total int) int {
	r := rand.New(rand.NewSource(seed))
	if total > 1000 {
		total = 1000
	}
	return r.Intn(total)
}