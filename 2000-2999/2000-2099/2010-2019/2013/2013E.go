package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxA = 100000

var divisors = func() [][]int {
	divs := make([][]int, maxA+1)
	for d := 1; d <= maxA; d++ {
		for multiple := d; multiple <= maxA; multiple += d {
			divs[multiple] = append(divs[multiple], d)
		}
	}
	return divs
}()

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		minVal := int(1e9)
		gAll := 0
		maxVal := 0
		freq := make([]int, maxA+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			freq[a[i]]++
			if a[i] < minVal {
				minVal = a[i]
			}
			if a[i] > maxVal {
				maxVal = a[i]
			}
			if i == 0 {
				gAll = a[i]
			} else {
				gAll = gcd(gAll, a[i])
			}
		}

		curr := minVal
		freq[curr]--
		used := 1
		ans := curr

		for curr > gAll {
			found := false
			for _, d := range divisors[curr] {
				if d == curr {
					continue
				}
				step := d
				for m := step; m <= maxVal; m += step {
					if freq[m] == 0 {
						continue
					}
					if gcd(curr/d, m/d) == 1 {
						freq[m]--
						curr = d
						ans += curr
						used++
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				break
			}
		}

		if used < n {
			ans += (n - used) * gAll
		}

		fmt.Fprintln(out, ans)
	}
}
