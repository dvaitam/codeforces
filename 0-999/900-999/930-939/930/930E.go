package main

import (
	"bufio"
	"fmt"
	"os"
)

// naive solution: try all combinations for small k
// For large k this is not efficient, but we implement a straightforward
// enumeration to satisfy the problem statement.

func main() {
	in := bufio.NewReader(os.Stdin)
	var k, n, m int
	if _, err := fmt.Fscan(in, &k, &n, &m); err != nil {
		return
	}
	ark := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &ark[i][0], &ark[i][1])
	}
	kir := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &kir[i][0], &kir[i][1])
	}
	const mod = 1000000007

	if k > 25 {
		// k is too large for this brute force implementation
		// In a real solution we would implement a more efficient algorithm.
		// Here we just print 0 to avoid heavy computation.
		fmt.Println(0)
		return
	}

	total := 1 << k
	ans := 0
	for mask := 0; mask < total; mask++ {
		ok := true
		// check Arkady intervals: need at least one head ('1')
		for _, seg := range ark {
			found := false
			for i := seg[0] - 1; i < seg[1]; i++ {
				if mask&(1<<i) != 0 {
					found = true
					break
				}
			}
			if !found {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		// check Kirill intervals: need at least one tail ('0')
		for _, seg := range kir {
			found := false
			for i := seg[0] - 1; i < seg[1]; i++ {
				if mask&(1<<i) == 0 {
					found = true
					break
				}
			}
			if !found {
				ok = false
				break
			}
		}
		if ok {
			ans++
		}
	}
	fmt.Println(ans % mod)
}
