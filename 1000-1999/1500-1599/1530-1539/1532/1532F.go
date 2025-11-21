package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	total := 2*n - 2
	strings := make([]string, total)
	lengthIdx := make([][]int, n)
	var longestIdx []int

	for i := 0; i < total; i++ {
		fmt.Fscan(in, &strings[i])
		l := len(strings[i])
		lengthIdx[l] = append(lengthIdx[l], i)
		if l == n-1 {
			longestIdx = append(longestIdx, i)
		}
	}

	if len(longestIdx) != 2 {
		return
	}

	s1 := strings[longestIdx[0]]
	s2 := strings[longestIdx[1]]

	candidates := []string{
		s1 + s2[len(s2)-1:],
		s2 + s1[len(s1)-1:],
	}

	for _, cand := range candidates {
		if len(cand) != n {
			continue
		}
		if ok, res := tryAssign(cand, strings, lengthIdx, n); ok {
			fmt.Fprintln(out, string(res))
			return
		}
	}
}

func tryAssign(cand string, strings []string, lengthIdx [][]int, n int) (bool, []byte) {
	total := len(strings)
	result := make([]byte, total)

	for length := 1; length < n; length++ {
		idxs := lengthIdx[length]
		if len(idxs) != 2 {
			return false, nil
		}
		prefix := cand[:length]
		suffix := cand[n-length:]

		i1, i2 := idxs[0], idxs[1]
		s1, s2 := strings[i1], strings[i2]

		if s1 == prefix && s2 == suffix {
			result[i1] = 'P'
			result[i2] = 'S'
			continue
		}
		if s1 == suffix && s2 == prefix {
			result[i1] = 'S'
			result[i2] = 'P'
			continue
		}
		return false, nil
	}

	return true, result
}
