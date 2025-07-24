package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func simulateStep(pos []int64, mov []bool, d int64) ([]int64, []bool) {
	n := len(pos)
	idx := 0
	for idx < n && !mov[idx] {
		idx++
	}
	if idx == n {
		return pos, mov
	}
	start := pos[idx]
	pos = append(pos[:idx], pos[idx+1:]...)
	flag := mov[idx]
	mov = append(mov[:idx], mov[idx+1:]...)
	cur := start
	e := d
	j := idx
	for e > 0 {
		if j >= len(pos) {
			cur += e
			break
		}
		gap := pos[j] - cur - 1
		if e > gap {
			e -= gap
			cur = pos[j]
			j++
		} else {
			cur += e
			e = 0
		}
	}
	insert := sort.Search(len(pos), func(i int) bool { return pos[i] >= cur })
	pos = append(pos, 0)
	copy(pos[insert+1:], pos[insert:])
	pos[insert] = cur
	mov = append(mov, false)
	copy(mov[insert+1:], mov[insert:])
	mov[insert] = flag
	return pos, mov
}

func solveQuery(n int, d int64, a []int64, s string, k int64, m int) int64 {
	pos := make([]int64, n)
	copy(pos, a)
	mov := make([]bool, n)
	for i := 0; i < n; i++ {
		mov[i] = s[i] == '1'
	}
	for t := int64(0); t < k; t++ {
		pos, mov = simulateStep(pos, mov, d)
	}
	return pos[m-1]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	var d int64
	var q int
	if _, err := fmt.Fscan(in, &n, &d, &q); err != nil {
		return
	}
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	var s string
	fmt.Fscan(in, &s)
	for ; q > 0; q-- {
		var k int64
		var m int
		fmt.Fscan(in, &k, &m)
		ans := solveQuery(n, d, a, s, k, m)
		fmt.Fprintln(out, ans)
	}
}
