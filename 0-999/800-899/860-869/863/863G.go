package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	m      int
	a      []int
	b      []int
	orient []int
	cur    []int
	path   [][]int
)

func gen(level int) {
	if level == 0 {
		s := make([]int, m)
		copy(s, cur)
		path = append(path, s)
		return
	}
	n := a[level-1]
	for i := 0; i < n; i++ {
		if orient[level-1] == 0 {
			cur[level-1] = i
		} else {
			cur[level-1] = n - 1 - i
		}
		gen(level - 1)
		if i != n-1 && level >= 2 {
			orient[level-2] ^= 1
		}
	}
}

func diffInstr(x, y []int) string {
	for i := 0; i < m; i++ {
		if x[i] != y[i] {
			if y[i] == x[i]+1 {
				return fmt.Sprintf("inc %d", i+1)
			}
			return fmt.Sprintf("dec %d", i+1)
		}
	}
	return ""
}

func main() {
	in := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(in, &m); err != nil {
		return
	}
	a = make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &a[i])
	}
	b = make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &b[i])
		b[i]--
	}
	orient = make([]int, m)
	cur = make([]int, m)
	gen(m)
	p := len(path)
	pos := 0
	for i := 0; i < p; i++ {
		ok := true
		for j := 0; j < m; j++ {
			if path[i][j] != b[j] {
				ok = false
				break
			}
		}
		if ok {
			pos = i
			break
		}
	}
	rotated := append(path[pos:], path[:pos]...)
	fmt.Println("Path")
	for i := 0; i < p-1; i++ {
		fmt.Println(diffInstr(rotated[i], rotated[i+1]))
	}
}
