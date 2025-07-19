package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func calc(x, pos int, opt []int) int {
	// iterative version of recursive calc
	sum := 0
	// opt is sorted in ascending order, length == original cnt
	for pos > 1 && x%2 == 0 && opt[pos-2] >= x/2 {
		sum += x
		x /= 2
		pos--
	}
	return sum + x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)
	// build frequency list
	var opt []int
	cnt := 1
	for i := 1; i < n; i++ {
		if a[i] != a[i-1] {
			opt = append(opt, cnt)
			cnt = 1
		} else {
			cnt++
		}
	}
	opt = append(opt, cnt)
	m := len(opt)
	if m == 1 {
		// only one distinct element
		fmt.Fprintln(writer, n)
		return
	}
	sort.Ints(opt)
	ans := opt[m-1]
	// make largest even
	if opt[m-1]%2 != 0 {
		opt[m-1]--
	}
	// try decreasing largest by 2
	for opt[m-1] > 0 {
		if ans > opt[m-1]*2 {
			break
		}
		p := calc(opt[m-1], m, opt)
		if p > ans {
			ans = p
		}
		opt[m-1] -= 2
	}
	fmt.Fprintln(writer, ans)
}
