package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	present := make(map[int]*list.Element)
	lst := list.New()

	for i := 0; i < n; i++ {
		var id int
		fmt.Fscan(in, &id)
		if _, ok := present[id]; ok {
			continue
		}
		if lst.Len() == k {
			backElem := lst.Back()
			backID := backElem.Value.(int)
			lst.Remove(backElem)
			delete(present, backID)
		}
		elem := lst.PushFront(id)
		present[id] = elem
	}

	fmt.Fprintln(out, lst.Len())
	first := true
	for e := lst.Front(); e != nil; e = e.Next() {
		if !first {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, e.Value.(int))
		first = false
	}
	if lst.Len() > 0 {
		fmt.Fprintln(out)
	}
}
