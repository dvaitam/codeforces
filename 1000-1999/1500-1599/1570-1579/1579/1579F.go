package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, d int
		fmt.Fscan(reader, &n, &d)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		g := gcd(n, d)
		ans := 0
		impossible := false
		for start := 0; start < g && !impossible; start++ {
			// collect elements of this cycle
			cycle := []int{}
			j := start
			for {
				cycle = append(cycle, arr[j])
				j = (j + d) % n
				if j == start {
					break
				}
			}
			allOne := true
			for _, v := range cycle {
				if v == 0 {
					allOne = false
					break
				}
			}
			if allOne {
				impossible = true
				break
			}
			// compute longest consecutive ones in circular cycle
			cur, best := 0, 0
			L := len(cycle)
			for i := 0; i < 2*L; i++ {
				if cycle[i%L] == 1 {
					cur++
					if cur > best {
						best = cur
					}
				} else {
					cur = 0
				}
			}
			if best > ans {
				ans = best
			}
		}
		if impossible {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}
