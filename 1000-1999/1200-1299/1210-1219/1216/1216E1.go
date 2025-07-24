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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}

	// Precompute prefix sums until total digits exceed 1e9
	const limit = 1000000000
	digitsPrefix := []int{0} // f[i] = sum_{j=1..i} len(j)
	blocksPrefix := []int{0} // g[i] = sum_{j=1..i} f[j]
	sumDigits := 0
	sumBlocks := 0
	for i := 1; sumBlocks < limit; i++ {
		// length of i in decimal
		l := 0
		for x := i; x > 0; x /= 10 {
			l++
		}
		sumDigits += l
		sumBlocks += sumDigits
		digitsPrefix = append(digitsPrefix, sumDigits)
		blocksPrefix = append(blocksPrefix, sumBlocks)
	}

	// helper to get digit at position idx (1-indexed) within concatenation of numbers 1..n
	digitInBlock := func(n int, idx int) byte {
		for d := 1; ; d++ {
			start := 1
			for i := 1; i < d; i++ {
				start *= 10
			}
			end := start*10 - 1
			if end > n {
				end = n
			}
			if end < start {
				continue
			}
			cnt := end - start + 1
			total := cnt * d
			if idx <= total {
				number := start + (idx-1)/d
				digitIdx := (idx - 1) % d
				// get digit at digitIdx
				digits := make([]byte, d)
				for i := d - 1; i >= 0; i-- {
					digits[i] = byte('0' + number%10)
					number /= 10
				}
				return digits[digitIdx]
			}
			idx -= total
		}
	}

	for ; q > 0; q-- {
		var k int
		fmt.Fscan(in, &k)
		// find block
		l, r := 1, len(blocksPrefix)-1
		for l < r {
			m := (l + r) / 2
			if blocksPrefix[m] < k {
				l = m + 1
			} else {
				r = m
			}
		}
		block := l
		prev := blocksPrefix[block-1]
		idx := k - prev
		// get digit
		b := digitInBlock(block, idx)
		fmt.Fprintln(out, string([]byte{b}))
	}
}
