package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxBit = 30 // enough for values up to (and above) 2^29

// nextVal returns the smallest value >= lower that has all bits in ones set
// and all bits in zeros unset. If impossible, returns -1.
func nextVal(lower, ones, zeros int64) int64 {
	if ones&zeros != 0 {
		return -1
	}

	// dp[pos][eq] = feasibility where pos is current bit, eq=1 if prefix==lower so far.
	var dp [maxBit + 1][2]int8
	for i := 0; i <= maxBit; i++ {
		for j := 0; j < 2; j++ {
			dp[i][j] = -1
		}
	}

	var feasible func(pos int, eq int) bool
	feasible = func(pos int, eq int) bool {
		if pos < 0 {
			return true
		}
		if dp[pos][eq] != -1 {
			return dp[pos][eq] == 1
		}
		bitL := int((lower >> pos) & 1)

		var opts [2]int
		cnt := 0
		if (ones>>pos)&1 == 1 {
			opts[0] = 1
			cnt = 1
		} else if (zeros>>pos)&1 == 1 {
			opts[0] = 0
			cnt = 1
		} else {
			opts[0] = 0
			opts[1] = 1
			cnt = 2
		}

		for i := 0; i < cnt; i++ {
			bit := opts[i]
			if eq == 1 && bit < bitL {
				continue
			}
			nextEq := eq
			if eq == 1 && bit > bitL {
				nextEq = 0
			}
			if feasible(pos-1, nextEq) {
				dp[pos][eq] = 1
				return true
			}
		}
		dp[pos][eq] = 0
		return false
	}

	if !feasible(maxBit, 1) {
		return -1
	}

	// Reconstruct minimal value.
	var res int64
	eq := 1
	for pos := maxBit; pos >= 0; pos-- {
		bitL := int((lower >> pos) & 1)

		var opts [2]int
		cnt := 0
		if (ones>>pos)&1 == 1 {
			opts[0] = 1
			cnt = 1
		} else if (zeros>>pos)&1 == 1 {
			opts[0] = 0
			cnt = 1
		} else {
			opts[0] = 0
			opts[1] = 1
			cnt = 2
		}

		chosen := 0
		for i := 0; i < cnt; i++ {
			bit := opts[i]
			if eq == 1 && bit < bitL {
				continue
			}
			nextEq := eq
			if eq == 1 && bit > bitL {
				nextEq = 0
			}
			if feasible(pos-1, nextEq) {
				chosen = bit
				eq = nextEq
				break
			}
		}
		if chosen == 1 {
			res |= 1 << pos
		}
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)

		a := make([]int64, n)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}

		// required ones per position
		req := make([]int64, n)
		req[0] = a[0]
		req[n-1] = a[n-2]
		for i := 1; i < n-1; i++ {
			req[i] = a[i-1] | a[i]
		}

		var forcedZero int64 = 0
		var total int64 = 0
		ok := true
		for i := 0; i < n; i++ {
			if req[i]&forcedZero != 0 {
				ok = false
				break
			}
			val := nextVal(b[i], req[i], forcedZero)
			if val == -1 {
				ok = false
				break
			}
			total += val - b[i]
			if i < n-1 {
				forcedZero = val & ^a[i]
			}
		}

		if ok {
			fmt.Fprintln(out, total)
		} else {
			fmt.Fprintln(out, -1)
		}
	}
}
