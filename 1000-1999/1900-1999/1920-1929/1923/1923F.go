package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int64 = 1_000_000_007

func swapStrategy(orig []byte, k int) []byte {
	s := make([]byte, len(orig))
	copy(s, orig)
	ones := make([]int, 0)
	zeros := make([]int, 0)
	for i, c := range s {
		if c == '1' {
			ones = append(ones, i)
		} else {
			zeros = append(zeros, i)
		}
	}
	i := 0
	j := len(zeros) - 1
	for i < len(ones) && j >= 0 && k > 0 && ones[i] < zeros[j] {
		s[ones[i]], s[zeros[j]] = s[zeros[j]], s[ones[i]]
		i++
		j--
		k--
	}
	return s
}

func valueMod(s []byte) int64 {
	var res int64
	for _, c := range s {
		res = (res*2 + int64(c-'0')) % MOD
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(reader, &s)

	first1 := -1
	last1 := -1
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			if first1 == -1 {
				first1 = i
			}
			last1 = i
		}
	}
	if first1 == -1 {
		fmt.Println(0)
		return
	}

	best := int64(-1)
	for r := 0; r <= 2 && r <= k; r++ {
		var t []byte
		if r == 0 {
			t = []byte(s)
		} else if r == 1 {
			tmp := []byte(s[first1:])
			for i, j := 0, len(tmp)-1; i < j; i, j = i+1, j-1 {
				tmp[i], tmp[j] = tmp[j], tmp[i]
			}
			t = tmp
		} else { // r == 2
			if last1 < first1 {
				continue
			}
			tmp := []byte(s[first1 : last1+1])
			t = tmp
		}
		t = swapStrategy(t, k-r)
		val := valueMod(t)
		if best == -1 || val < best {
			best = val
		}
	}
	fmt.Println(best)
}
