package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)

	weights := make([]int64, 10)
	leading := make([]bool, 10)
	pow10 := [7]int64{1}
	for i := 1; i < 7; i++ {
		pow10[i] = pow10[i-1] * 10
	}

	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		leading[s[0]-'a'] = true
		L := len(s)
		for j := L - 1; j >= 0; j-- {
			letter := s[j] - 'a'
			pos := L - 1 - j
			weights[letter] += pow10[pos]
		}
	}

	bestSum := int64(1<<62 - 1)
	for zero := 0; zero < 10; zero++ {
		if leading[zero] {
			continue
		}
		type pair struct {
			w   int64
			idx int
		}
		pairs := make([]pair, 0, 9)
		for i := 0; i < 10; i++ {
			if i == zero {
				continue
			}
			pairs = append(pairs, pair{weights[i], i})
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].w == pairs[j].w {
				return pairs[i].idx < pairs[j].idx
			}
			return pairs[i].w > pairs[j].w
		})
		digits := make([]int, 10)
		digits[zero] = 0
		for i, p := range pairs {
			digits[p.idx] = i + 1
		}
		var sum int64
		for i := 0; i < 10; i++ {
			sum += weights[i] * int64(digits[i])
		}
		if sum < bestSum {
			bestSum = sum
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, bestSum)
	out.Flush()
}
