package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)
	n := len(s)
	const AL = 20
	M := 1 << AL
	const INF = int(1e9)
	earliest := make([]int, M)
	latest := make([]int, M)
	for i := 0; i < M; i++ {
		earliest[i] = INF
		latest[i] = -1
	}
	ans := 0
	bs := []byte(s)
	for i := 0; i < n; i++ {
		mask := 0
		for j := i; j < n && j-i < AL; j++ {
			c := int(bs[j] - 'a')
			if mask>>c&1 == 1 {
				break
			}
			mask |= 1 << c
			l := j - i + 1
			if earliest[mask] > j {
				earliest[mask] = j
			}
			if latest[mask] < i {
				latest[mask] = i
			}
			if l > ans {
				ans = l
			}
		}
	}

	full := M - 1
	for m1 := 1; m1 < M; m1++ {
		if earliest[m1] == INF {
			continue
		}
		comp := full ^ m1
		sub := comp
		for sub > 0 {
			if latest[sub] > earliest[m1] {
				l := bitsOnesCount(uint(sub)) + bitsOnesCount(uint(m1))
				if l > ans {
					ans = l
				}
			}
			sub = (sub - 1) & comp
		}
	}
	fmt.Println(ans)
}

func bitsOnesCount(x uint) int {
	return int(bits.OnesCount(x))
}
