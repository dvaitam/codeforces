package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Ints(a)
	prefix := make([]int64, n+1)
	for i, v := range a {
		prefix[i+1] = prefix[i] + int64(v)
	}
	ans := int64(1 << 62)
	i := 0
	for i < n {
		j := i
		for j < n && a[j] == a[i] {
			j++
		}
		x := a[i]
		cnt := j - i
		if cnt >= k {
			fmt.Fprintln(out, 0)
			return
		}
		need := k - cnt
		left := i
		right := n - j
		if left >= need {
			cost := int64(x)*int64(need) - (prefix[left] - prefix[left-need])
			if cost < ans {
				ans = cost
			}
		}
		if right >= need {
			cost := (prefix[j+need] - prefix[j]) - int64(x)*int64(need)
			if cost < ans {
				ans = cost
			}
		}
		if left+right >= need {
			if left < need {
				leftCost := int64(x)*int64(left) - prefix[left]
				rightTake := need - left
				rightCost := (prefix[j+rightTake] - prefix[j]) - int64(x)*int64(rightTake)
				if leftCost+rightCost < ans {
					ans = leftCost + rightCost
				}
			}
			if right < need {
				rightCost := prefix[n] - prefix[j] - int64(x)*int64(right)
				leftTake := need - right
				leftCost := int64(x)*int64(leftTake) - (prefix[left] - prefix[left-leftTake])
				if leftCost+rightCost < ans {
					ans = leftCost + rightCost
				}
			}
		}
		i = j
	}
	if ans < 0 {
		ans = 0
	}
	fmt.Fprintln(out, ans)
}
