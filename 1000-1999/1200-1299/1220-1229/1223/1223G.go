package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func chk(x int64, y int, ans *int64) {
	if x > 1 {
		area := x * int64(y)
		if area > *ans {
			*ans = area
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	maxLen := 0
	const N = 1000005
	a := make([]int, N)
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(in, &x)
		if x >= len(a) {
			continue
		}
		a[x]++
		if x > maxLen {
			maxLen = x
		}
	}
	l := make([]int, N)
	l[0] = -1
	for i := 1; i <= maxLen*2 && i < N; i++ {
		if a[i] > 0 {
			l[i] = i
		} else {
			l[i] = l[i-1]
		}
		a[i] += a[i-1]
	}

	ans := int64(0)
	for y := 2; y <= maxLen; y++ {
		cnt := int64(0)
		for i := y - 1; i <= maxLen; i += y {
			if i+y >= N {
				break
			}
			cnt += int64(i/y+1) * int64(a[i+y]-a[i])
		}
		mx1, mx2 := -1, -1
		j := (maxLen / y) * y
		for i := maxLen; i >= 0; i = j - 1 {
			t1 := l[i]
			var t2 int
			if t1 >= 0 && a[t1]-a[t1-1] > 1 {
				t2 = t1
			} else {
				if t1 < 0 {
					t2 = -1
				} else {
					t2 = l[t1-1]
				}
			}
			k := j / y
			if t1 >= j {
				t1mod := t1 % y
				if mx1 != -1 && (mx2 != -1 || t1mod >= mx1) {
					cand := min64(cnt-int64(k*2)-1, int64(j+max(t1mod, mx1)))
					chk(cand, y, &ans)
				}
				if t1mod > mx1 {
					mx2 = mx1
					mx1 = t1mod
				} else if t1mod > mx2 {
					mx2 = t1mod
				}
				if t2 >= j {
					if r := t2 % y; r > mx2 {
						mx2 = r
					}
				}
			}
			if mx1 != -1 {
				cand1 := min64(cnt-int64(k), int64(j+mx1)/2)
				chk(cand1, y, &ans)
				if mx2 != -1 {
					cand2 := min64(cnt-int64(k*2), int64(j+mx2))
					chk(cand2, y, &ans)
				}
			}
			if j == 0 {
				break
			}
			j -= y
		}
	}

	fmt.Println(ans)
}
