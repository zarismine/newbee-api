package passwd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenOrderNo() string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < 4; i++ {
		_, err := fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
		if err != nil {
			return ""
		}
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	return timestamp + sb.String()
}