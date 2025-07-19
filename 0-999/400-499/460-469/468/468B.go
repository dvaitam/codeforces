package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}
	sw := 0
	if b < a {
		a, b = b, a
		sw = 1
	}
	arr := make([]int, n)
	m := make(map[int]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		m[arr[i]] = i
	}
	bi := make([]int, n)
	for _, x := range arr {
		idx, ok := m[x]
		if !ok {
			continue
		}
		// try pairing for sum b
		if j, ok2 := m[b-x]; ok2 {
			s := 1 - sw
			bi[idx] = s
			bi[j] = s
			delete(m, x)
			delete(m, b-x)
		} else if j, ok2 := m[a-x]; ok2 {
			s := sw
			bi[idx] = s
			bi[j] = s
			delete(m, x)
			delete(m, a-x)
		} else {
			fmt.Fprintln(out, "NO")
			return
		}
	}
	fmt.Fprintln(out, "YES")
	for i := 0; i < n; i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprint(bi[i]))
	}
	out.WriteByte('\n')
}
