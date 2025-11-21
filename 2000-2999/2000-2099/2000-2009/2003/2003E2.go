package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf = int64(-1 << 60)

func comb(x int64) int64 {
	return x * (x - 1) / 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)

	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)

		needZero := make([]bool, n+2)
		needOne := make([]bool, n+2)
		coverDiff := make([]int, n+3)

		for i := 0; i < m; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			needZero[l] = true
			needOne[r] = true
			coverDiff[l]++
			coverDiff[r]--
		}

		impossible := false
		for i := 1; i <= n; i++ {
			if needZero[i] && needOne[i] {
				impossible = true
				break
			}
		}
		if impossible {
			fmt.Fprintln(out, -1)
			continue
		}

		cover := make([]bool, n+1)
		cur := 0
		for i := 1; i <= n-1; i++ {
			cur += coverDiff[i]
			if cur > 0 {
				cover[i] = true
			}
		}

		type component struct {
			length  int
			minZero int
			maxZero int
		}
		components := make([]component, 0)
		for i := 1; i <= n; {
			j := i
			for j < n && cover[j] {
				j++
			}
			length := j - i + 1
			minZero := 0
			maxZero := length

			latestZero := 0
			earliestOne := length + 1
			for pos := i; pos <= j; pos++ {
				offset := pos - i + 1 // positions numbered starting from 1
				if needZero[pos] {
					if offset > latestZero {
						latestZero = offset
					}
				}
				if needOne[pos] {
					if offset < earliestOne {
						earliestOne = offset
					}
				}
			}
			minZero = latestZero
			maxZero = earliestOne - 1
			if earliestOne == length+1 {
				maxZero = length
			}

			if minZero > maxZero {
				impossible = true
				break
			}
			components = append(components, component{length: length, minZero: minZero, maxZero: maxZero})
			i = j + 1
		}

		if impossible {
			fmt.Fprintln(out, -1)
			continue
		}

		dp := make([]int64, n+1)
		for i := range dp {
			dp[i] = negInf
		}
		dp[0] = 0

		for _, comp := range components {
			next := make([]int64, n+1)
			for i := range next {
				next[i] = negInf
			}

			for onesBefore := 0; onesBefore <= n; onesBefore++ {
				if dp[onesBefore] == negInf {
					continue
				}
				for zeros := comp.minZero; zeros <= comp.maxZero; zeros++ {
					if zeros < 0 || zeros > comp.length {
						continue
					}
					ones := comp.length - zeros
					val := dp[onesBefore] + int64(onesBefore)*int64(zeros)
					newOnes := onesBefore + ones
					if val > next[newOnes] {
						next[newOnes] = val
					}
				}
			}
			dp = next
		}

		ans := int64(negInf)
		for onesCount, v := range dp {
			if v == negInf {
				continue
			}
			total := v + comb(int64(onesCount)) + comb(int64(n-onesCount))
			if total > ans {
				ans = total
			}
		}
		fmt.Fprintln(out, ans)
	}
}
