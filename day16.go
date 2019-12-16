package main

import (
	"log"
	"strings"
)

func day16() {
	src := strings.TrimSpace(mustdaydatastr(16))

	log.Println("day16a:", phasesigN(src, 100)[:8])
}

func phasesigN(sig string, n int) string {
	sigb := []byte(sig)
	for i := 0; i < n; i++ {
		phasesigb(sigb)
	}
	return string(sigb)
}

func phasesigb(sig []byte) {
	src := make([]int, 0, len(sig))
	for _, b := range sig {
		src = append(src, int(b-'0'))
	}

	for i := range sig {
		w := 0
		for j, v := range src {
			w += v * repat(i, j)
		}
		if w < 0 {
			w = -w
		}
		sig[i] = '0' + byte(w%10)
	}
}

func repat(rpt int, p int) int {
	values := []int{0, 1, 0, -1}
	i := ((p + 1) / (rpt + 1)) % len(values)
	return values[i]
}
