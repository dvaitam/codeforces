package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxM = 5000000

var spf [maxM + 1]int

func init() {
	for i := 2; i <= maxM; i++ {
		if spf[i] == 0 {
			spf[i] = i
			if i <= maxM/i {
				for j := i * i; j <= maxM; j += i {
					if spf[j] == 0 {
						spf[j] = i
					}
				}
			}
		}
	}
	spf[1] = 1
}

func maxPrimeFactor(x int) int {
	res := 1
	for x > 1 {
		p := spf[x]
		if p > res {
			res = p
		}
		for x%p == 0 {
			x /= p
		}
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		arr := make([]int, n)
		maxVal := 0
		maxPF := 1
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
			if arr[i] > maxVal {
				maxVal = arr[i]
			}
			pf := maxPrimeFactor(arr[i])
			if pf > maxPF {
				maxPF = pf
			}
		}
		if maxVal == 1 {
			fmt.Fprintln(writer, 0)
			continue
		}
		freq := make([]int, maxVal+1)
		for _, v := range arr {
			freq[v]++
		}
		dp := make([]int, maxVal+1)
		cnt := make([]int, maxVal+1)
		for val, c := range freq {
			if c > 0 {
				if val == 1 {
					dp[val] = 1
				}
				cnt[dp[val]] += c
			}
		}
		curMin := 0
		for curMin < len(cnt) && cnt[curMin] == 0 {
			curMin++
		}
		best := int(1 << 60)
		for r := 2; r <= maxVal; r++ {
			if r > dp[r] {
				old := dp[r]
				dp[r] = r
				if freq[r] > 0 {
					cnt[old] -= freq[r]
					cnt[r] += freq[r]
					if old == curMin && cnt[old] == 0 {
						for curMin < len(cnt) && cnt[curMin] == 0 {
							curMin++
						}
					}
				}
			}
			for j := r * 2; j <= maxVal; j += r {
				cand := dp[j/r]
				if cand == 0 {
					continue
				}
				if cand > r {
					cand = r
				}
				if cand > dp[j] {
					old := dp[j]
					dp[j] = cand
					if freq[j] > 0 {
						cnt[old] -= freq[j]
						cnt[cand] += freq[j]
						if old == curMin && cnt[old] == 0 {
							for curMin < len(cnt) && cnt[curMin] == 0 {
								curMin++
							}
						}
					}
				}
			}
			if r >= maxPF {
				if curMin < len(cnt) {
					diff := r - curMin
					if diff < best {
						best = diff
						if best == 0 {
							break
						}
					}
				}
			}
		}
		if best == int(1<<60) {
			best = 0
		}
		fmt.Fprintln(writer, best)
	}
}
