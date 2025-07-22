package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	base1 uint64 = 911382323
	base2 uint64 = 972663749
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var N, M int
	if _, err := fmt.Fscan(in, &N, &M); err != nil {
		return
	}
	first := make([][]byte, N)
	for i := 0; i < N; i++ {
		var s string
		fmt.Fscan(in, &s)
		first[i] = []byte(s)
	}
	second := make([][]byte, M)
	for i := 0; i < M; i++ {
		var s string
		fmt.Fscan(in, &s)
		second[i] = []byte(s)
	}

	pow1 := make([]uint64, N+1)
	pow2 := make([]uint64, N+1)
	pow1[0], pow2[0] = 1, 1
	for i := 1; i <= N; i++ {
		pow1[i] = pow1[i-1] * base1
		pow2[i] = pow2[i-1] * base2
	}

	rowHash1 := make([]uint64, N)
	for i := 0; i < N; i++ {
		var h uint64
		for j := 0; j < M; j++ {
			h = h*base1 + uint64(first[i][j])
		}
		rowHash1[i] = h
	}

	prefixV := make([]uint64, N+1)
	for i := 0; i < N; i++ {
		prefixV[i+1] = prefixV[i]*base2 + rowHash1[i]
	}
	vertical1 := make([]uint64, N-M+1)
	pow2M := pow2[M]
	for i := 0; i <= N-M; i++ {
		vertical1[i] = prefixV[i+M] - prefixV[i]*pow2M
	}

	rowHash2 := make([][]uint64, M)
	pow1M := pow1[M]
	for r := 0; r < M; r++ {
		prefix := make([]uint64, N+1)
		for c := 0; c < N; c++ {
			prefix[c+1] = prefix[c]*base1 + uint64(second[r][c])
		}
		rowHash2[r] = make([]uint64, N-M+1)
		for j := 0; j <= N-M; j++ {
			rowHash2[r][j] = prefix[j+M] - prefix[j]*pow1M
		}
	}

	finalHash2 := make([]uint64, N-M+1)
	for j := 0; j <= N-M; j++ {
		var h uint64
		for r := 0; r < M; r++ {
			h = h*base2 + rowHash2[r][j]
		}
		finalHash2[j] = h
	}

	mp := make(map[uint64][]int)
	for j := 0; j <= N-M; j++ {
		mp[finalHash2[j]] = append(mp[finalHash2[j]], j)
	}

	for i := 0; i <= N-M; i++ {
		if list, ok := mp[vertical1[i]]; ok {
			for _, j := range list {
				if checkEqual(first, second, i, j, M) {
					fmt.Fprintf(out, "%d %d\n", i+1, j+1)
					return
				}
			}
		}
	}
}

func checkEqual(first [][]byte, second [][]byte, i, j, M int) bool {
	for r := 0; r < M; r++ {
		for c := 0; c < M; c++ {
			if first[i+r][c] != second[r][j+c] {
				return false
			}
		}
	}
	return true
}
