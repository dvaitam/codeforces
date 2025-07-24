package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	ink := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &ink[i])
	}

	L := int64(n) * 7 / gcd(int64(n), 7)
	pos := make([][]int64, n)
	for d := int64(1); d <= L; d++ {
		pen := int((d - 1) % int64(n))
		dow := (d - 1) % 7
		if dow != 6 {
			pos[pen] = append(pos[pen], d)
		}
	}

	bestPen := -1
	bestDay := int64(1<<63 - 1)
	for i := 0; i < n; i++ {
		arr := pos[i]
		cnt := int64(len(arr))
		if cnt == 0 {
			continue
		}
		k := ink[i]
		q := (k - 1) / cnt
		p := arr[(k-1)%cnt]
		day := q*L + p
		if day < bestDay {
			bestDay = day
			bestPen = i + 1
		}
	}

	fmt.Println(bestPen)
}
