package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

func solve(r *bufio.Reader, w *bufio.Writer) {
	var n, m int
	fmt.Fscan(r, &n, &m)

	actions := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &actions[i])
	}

	ans := make([]int, n)
	for i := range ans {
		ans[i] = -1
	}

	// list from top to bottom: front is top
	l := list.New()
	pos := make(map[int]*list.Element, n+m)

	for i := 1; i <= n; i++ {
		el := l.PushBack(i)
		pos[i] = el
	}

	for i, v := range actions {
		if el, ok := pos[v]; ok {
			l.MoveToFront(el)
		} else {
			// remove bottom element
			back := l.Back()
			removed := back.Value.(int)
			l.Remove(back)
			delete(pos, removed)
			if removed <= n && ans[removed-1] == -1 {
				ans[removed-1] = i + 1
			}
			pos[v] = l.PushFront(v)
		}
	}

	for i, v := range ans {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, v)
	}
	fmt.Fprintln(w)
}
