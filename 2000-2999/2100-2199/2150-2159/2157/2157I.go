package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const M = 1000010

var alive = make([]int, M)
var whenArr = make([]int, M)
var lose = make([][]int, M)

func main() {
	for i := 0; i < M; i++ {
		alive[i] = -1
		whenArr[i] = -1
	}

	for m := 3; m < M; m += 2 {
		bases := []int{0, m + 1}
		base := m + 1
		lose[m] = append(lose[m], base)

		for base+(m+1) < M {
			var dfs func(int) int
			dfs = func(i int) int {
				if i < m {
					return 0
				}
				if i == m {
					return 1
				}
				if whenArr[i] == m {
					return alive[i]
				}
				whenArr[i] = m
				pos := sort.SearchInts(bases, i)
				pos--
				myBase := bases[pos]
				d := i - myBase
				if d == m+1 {
					if dfs(myBase+d/2) == 1 {
						alive[i] = 1
					} else {
						alive[i] = 0
					}
					return alive[i]
				}
				holes := 1
				if d%2 == 0 {
					if dfs(myBase+d/2) == 1 {
						holes--
						if holes == 0 {
							alive[i] = 0
							return 0
						}
					}
				}
				prevBase := bases[pos-1]
				diffPrev := myBase - prevBase
				if diffPrev == m+1 {
					if (d+m+1)%2 == 0 {
						if dfs(prevBase+(d+m+1)/2) == 1 {
							holes--
							if holes == 0 {
								alive[i] = 0
								return 0
							}
						}
					}
				} else {
					if d == m/2 {
						holes--
						if holes == 0 {
							alive[i] = 0
							return 0
						}
					}
					if d < m && (d+m+2)%2 == 0 {
						if dfs(prevBase+(d+m+2)/2) == 1 {
							holes--
							if holes == 0 {
								alive[i] = 0
								return 0
							}
						}
					}
				}
				alive[i] = 1
				return 1
			}

			me := dfs(base + (m + 1))
			if me == 1 {
				base += m + 2
			} else {
				base += m + 1
			}
			lose[m] = append(lose[m], base)
			bases = append(bases, base)
		}
	}

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var tt int
	fmt.Fscan(in, &tt)
	for ; tt > 0; tt-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		if m%2 == 0 {
			if n%(m+1) == 0 {
				fmt.Fprintln(out, "NO")
			} else {
				fmt.Fprintln(out, "YES")
			}
		} else {
			pos := sort.SearchInts(lose[m], n)
			if pos < len(lose[m]) && lose[m][pos] == n {
				fmt.Fprintln(out, "NO")
			} else {
				fmt.Fprintln(out, "YES")
			}
		}
	}
}

