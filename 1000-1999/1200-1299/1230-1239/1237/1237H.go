package main

import (
	"bufio"
	"fmt"
	"os"
)

func sw(x int) int {
	if x == 1 {
		return 2
	}
	if x == 2 {
		return 1
	}
	return x
}

func pairCode(a, b byte) int {
	if a == '0' {
		if b == '0' {
			return 0
		}
		return 1
	}
	if b == '0' {
		return 2
	}
	return 3
}

func parsePairs(s string) []int {
	m := len(s) / 2
	res := make([]int, m)
	for i := 0; i < m; i++ {
		res[i] = pairCode(s[2*i], s[2*i+1])
	}
	return res
}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b string
		fmt.Fscan(in, &a, &b)

		n := len(a)
		m := n / 2

		pa := parsePairs(a)
		pb := parsePairs(b)

		ca := [3]int{}
		cb := [3]int{}
		for i := 0; i < m; i++ {
			switch pa[i] {
			case 0:
				ca[0]++
			case 3:
				ca[2]++
			default:
				ca[1]++
			}
			switch pb[i] {
			case 0:
				cb[0]++
			case 3:
				cb[2]++
			default:
				cb[1]++
			}
		}

		if ca != cb {
			fmt.Fprintln(out, -1)
			continue
		}

		ans := make([]int, 0, n+1)

		flip := func(r int) {
			for l, rr := 0, r-1; l < rr; l, rr = l+1, rr-1 {
				pa[l], pa[rr] = sw(pa[rr]), sw(pa[l])
			}
			if r%2 == 1 {
				pa[r/2] = sw(pa[r/2])
			}
			ans = append(ans, 2*r)
		}

		ok := true

		for i := m; i >= 1 && ok; i-- {
			tgt := pb[i-1]

			if tgt == 0 || tgt == 3 {
				k := 0
				for j := i; j >= 1; j-- {
					if pa[j-1] == tgt {
						k = j
						break
					}
				}
				if k == 0 {
					ok = false
					break
				}
				if k < i {
					if k == 1 {
						flip(i)
					} else {
						flip(k)
						flip(i)
					}
				}
			} else {
				k := 0
				for j := i; j >= 1; j-- {
					if pa[j-1] == tgt {
						k = j
						break
					}
				}
				if k != 0 {
					if k < i {
						flip(k)
						flip(i)
					}
				} else {
					opp := sw(tgt)
					k = 0
					for j := 1; j <= i; j++ {
						if pa[j-1] == opp {
							k = j
							break
						}
					}
					if k == 0 {
						ok = false
						break
					}
					if k == 1 {
						flip(i)
					} else {
						flip(k)
						flip(1)
						flip(i)
					}
				}
			}
		}

		if ok {
			for i := 0; i < m; i++ {
				if pa[i] != pb[i] {
					ok = false
					break
				}
			}
		}

		if !ok || len(ans) > n+1 {
			fmt.Fprintln(out, -1)
			continue
		}

		fmt.Fprint(out, len(ans))
		for _, x := range ans {
			fmt.Fprint(out, " ", x)
		}
		fmt.Fprintln(out)
	}
}