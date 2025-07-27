package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution for problemB.txt (Orac and Medians).
// We can choose any subsegment and replace all numbers in it with the
// median of that segment. The median is always one of the numbers from
// the segment, so the value k must already appear in the array to be
// attainable. Furthermore the operation can propagate a value at least k
// across neighbours when there are at least two numbers >= k within
// distance two. Therefore a configuration is transformable if and only
// if there exists an element equal to k and either n == 1 or there exist
// two elements not farther than two positions apart that are >= k.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		hasK := false
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] == k {
				hasK = true
			}
		}
		good := false
		for i := 0; i < n-1; i++ {
			if a[i] >= k && a[i+1] >= k {
				good = true
				break
			}
		}
		if !good {
			for i := 0; i < n-2; i++ {
				if a[i] >= k && a[i+2] >= k {
					good = true
					break
				}
			}
		}
		if hasK && (good || n == 1) {
			fmt.Fprintln(out, "yes")
		} else {
			fmt.Fprintln(out, "no")
		}
	}
}
