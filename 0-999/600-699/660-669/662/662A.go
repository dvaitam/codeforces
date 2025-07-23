package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	base := uint64(0)
	diffs := make([]uint64, n)
	for i := 0; i < n; i++ {
		var a, b uint64
		fmt.Fscan(reader, &a, &b)
		base ^= a
		diffs[i] = a ^ b
	}

	var basis [64]uint64
	rank := 0
	for _, d := range diffs {
		x := d
		for j := 63; j >= 0; j-- {
			if (x>>uint(j))&1 == 0 {
				continue
			}
			if basis[j] == 0 {
				basis[j] = x
				rank++
				break
			}
			x ^= basis[j]
		}
	}

	// check if base is in span of diffs
	inSpan := func(v uint64) bool {
		x := v
		for j := 63; j >= 0; j-- {
			if (x>>uint(j))&1 == 0 {
				continue
			}
			if basis[j] == 0 {
				return false
			}
			x ^= basis[j]
		}
		return true
	}(base)

	if !inSpan {
		fmt.Fprint(writer, "1/1")
		return
	}

	denom := uint64(1) << uint(rank)
	numer := denom - 1
	fmt.Fprintf(writer, "%d/%d", numer, denom)
}
