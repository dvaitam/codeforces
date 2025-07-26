package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const MAXK = 2000

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	diff := make([]int, MAXK+2)
	for _, x := range a {
		sMin := int(math.Sqrt(float64(max(0, x-MAXK)))) - 2
		if sMin < 0 {
			sMin = 0
		}
		sMax := int(math.Sqrt(float64(x+MAXK))) + 2
		for s := sMin; s <= sMax; s++ {
			sq := s * s
			l := sq - x
			r := sq + s - x
			if r < 0 || l > MAXK {
				continue
			}
			if l < 0 {
				l = 0
			}
			if r > MAXK {
				r = MAXK
			}
			diff[l]++
			diff[r+1]--
		}
	}

	cnt := 0
	for k := 0; k <= MAXK; k++ {
		cnt += diff[k]
		if cnt == n {
			fmt.Println(k)
			return
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
