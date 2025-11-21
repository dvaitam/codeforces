package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n%2 == 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		ans := construct(n)
		if ans == nil {
			fmt.Fprintln(out, -1)
			continue
		}
		var sb strings.Builder
		for i, v := range ans {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		fmt.Fprintln(out, sb.String())
	}
}

func construct(n int) []int {
	ans := make([]int, n)
	used := make([]bool, n+1)
	label := 1
	for i := n; i >= 1; i-- {
		if used[i] {
			continue
		}
		found := false
		maxS := int(math.Sqrt(float64(i - 1)))
		for s := maxS; s >= 1; s-- {
			j := i - s*s
			if j >= 1 && !used[j] {
				used[i], used[j] = true, true
				ans[i-1], ans[j-1] = label, label
				label++
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return ans
}
