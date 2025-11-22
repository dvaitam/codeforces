package main

import (
	"bufio"
	"fmt"
	"os"
)

// We pack the digit counts into a single integer using base 21 (counts never exceed 20).
// A vector can then be updated by adding basePowers[d] for each digit increment,
// without decoding.
const base = 21

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	basePow := [10]uint64{1}
	for i := 1; i < 10; i++ {
		basePow[i] = basePow[i-1] * base
	}

	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)

		digits := make([]int, 0, 12)
		tmp := n
		for tmp > 0 {
			digits = append(digits, int(tmp%10)) // least significant first
			tmp /= 10
		}
		if len(digits) == 0 {
			digits = append(digits, 0)
		}
		lenD := len(digits)

		curr := [4]map[uint64]struct{}{}
		next := [4]map[uint64]struct{}{}
		curr[1] = map[uint64]struct{}{0: {}}

		for pos := 0; pos < lenD; pos++ {
			limit := 9
			if digits[pos] >= 0 {
				limit = digits[pos]
			}
			for i := 0; i < 4; i++ {
				next[i] = nil
			}

			for state := 0; state < 4; state++ {
				if len(curr[state]) == 0 {
					continue
				}
				tight := state / 2
				carry := state % 2
				maxDigit := 9
				if tight == 1 {
					maxDigit = limit
				}

				for key := range curr[state] {
					for d := 0; d <= maxDigit; d++ {
						nextTight := 0
						if tight == 1 && d == limit {
							nextTight = 1
						}

						add := 0
						if pos == 0 {
							add = 1
						}
						sum := d + add + carry
						ndigit := sum
						nextCarry := 0
						if sum >= 10 {
							ndigit = sum - 10
							nextCarry = 1
						}

						nkey := key + basePow[d] + basePow[ndigit]
						if pos == lenD-1 && nextCarry == 1 {
							// extra highest digit 1
							nkey += basePow[1]
							nextCarry = 0
						}

						idx := nextTight*2 + nextCarry
						if next[idx] == nil {
							next[idx] = make(map[uint64]struct{})
						}
						next[idx][nkey] = struct{}{}
					}
				}
			}

			curr = next
		}

		unique := make(map[uint64]struct{})
		for state := 0; state < 4; state++ {
			for key := range curr[state] {
				unique[key] = struct{}{}
			}
		}

		fmt.Fprintln(out, len(unique))
	}
}
