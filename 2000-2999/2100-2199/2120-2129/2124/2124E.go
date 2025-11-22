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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		var sum, mx int64
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
			if a[i] > mx {
				mx = a[i]
			}
		}

		if sum%2 == 1 || mx > sum/2 {
			fmt.Fprintln(out, -1)
			continue
		}

		half := sum / 2
		pref := int64(0)
		equalIdx := -1
		for i := 0; i < n; i++ {
			pref += a[i]
			if pref == half {
				equalIdx = i
				break
			}
		}

		if equalIdx != -1 {
			fmt.Fprintln(out, 1)
			for i := 0; i < n; i++ {
				if i > 0 {
					fmt.Fprint(out, " ")
				}
				fmt.Fprint(out, a[i])
			}
			fmt.Fprintln(out)
			continue
		}

		// Two-operations construction.
		// Find first index where prefix exceeds half.
		pref = 0
		k := -1
		for i := 0; i < n; i++ {
			pref += a[i]
			if pref > half {
				k = i
				break
			}
		}

		prev := pref - a[k]      // prefix before k
		delta := half - prev     // positive and < a[k]
		b1 := make([]int64, n)   // first operation
		for i := 0; i < k; i++ { // zero out prefix before k
			b1[i] = a[i]
		}
		b1[k] = a[k] - delta
		// Remove remaining (half - a[k]) from suffix after k
		rem := half - a[k]
		for i := k + 1; i < n && rem > 0; i++ {
			if a[i] >= rem {
				b1[i] = rem
				rem = 0
			} else {
				b1[i] = a[i]
				rem -= a[i]
			}
		}

		aPrime := make([]int64, n)
		for i := 0; i < n; i++ {
			aPrime[i] = a[i] - b1[i]
		}
		// After first operation, total is 2*delta and prefix up to k equals delta,
		// so we can zero everything with b2 = aPrime.

		fmt.Fprintln(out, 2)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, b1[i])
		}
		fmt.Fprintln(out)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, aPrime[i])
		}
		fmt.Fprintln(out)
	}
}
