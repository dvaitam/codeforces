package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type pair struct {
	val   int
	count int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	freq := make(map[int]int)
	for _, v := range nums {
		freq[v]++
	}

	type freqPair struct{ val, freq int }
	arr := make([]freqPair, 0, len(freq))
	for v, f := range freq {
		arr = append(arr, freqPair{v, f})
	}

	sort.Slice(arr, func(i, j int) bool { return arr[i].freq > arr[j].freq })

	// gather frequencies only for quick calculations
	freqs := make([]int, len(arr))
	for i := range arr {
		freqs[i] = arr[i].freq
	}
	prefix := make([]int, len(freqs)+1)
	for i, f := range freqs {
		prefix[i+1] = prefix[i] + f
	}

	bestArea, bestP, bestQ := 0, 0, 0
	// iterate possible q (number of columns)
	for q := 1; q*q <= n; q++ {
		// find index where freq < q
		idx := sort.Search(len(freqs), func(i int) bool { return freqs[i] < q })
		total := prefix[idx] + (len(freqs)-idx)*q
		p := total / q
		if p >= q && p*q > bestArea {
			bestArea = p * q
			bestP = p
			bestQ = q
		}
	}

	// prepare values to fill matrix
	use := make([]pair, 0, len(arr))
	for _, fp := range arr {
		c := fp.freq
		if c > bestP {
			c = bestP
		}
		if c > bestQ {
			c = bestQ
		}
		if c > 0 {
			use = append(use, pair{fp.val, c})
		}
	}

	sort.Slice(use, func(i, j int) bool { return use[i].count > use[j].count })

	vals := make([]int, 0, bestArea)
	for _, p := range use {
		for k := 0; k < p.count && len(vals) < bestArea; k++ {
			vals = append(vals, p.val)
		}
	}

	// fill matrix diagonally
	matrix := make([][]int, bestP)
	for i := 0; i < bestP; i++ {
		matrix[i] = make([]int, bestQ)
	}
	idx := 0
	for col := 0; col < bestQ; col++ {
		row := 0
		for row < bestP && idx < bestArea {
			matrix[row][(row+col)%bestQ] = vals[idx]
			idx++
			row++
		}
	}

	fmt.Fprintln(out, bestArea)
	fmt.Fprintln(out, bestP, bestQ)
	for i := 0; i < bestP; i++ {
		for j := 0; j < bestQ; j++ {
			if j > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, matrix[i][j])
		}
		fmt.Fprintln(out)
	}
}
