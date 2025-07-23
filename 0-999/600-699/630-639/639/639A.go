package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct {
	id int
	t  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k, q int
	if _, err := fmt.Fscan(reader, &n, &k, &q); err != nil {
		return
	}
	tvals := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &tvals[i])
	}

	displayed := make([]pair, 0, k)
	inDisp := make([]bool, n+1)

	for ; q > 0; q-- {
		var typ, id int
		fmt.Fscan(reader, &typ, &id)
		if typ == 1 {
			if k == 0 {
				continue
			}
			// insert friend if belongs to top k
			if len(displayed) < k {
				displayed = append(displayed, pair{id: id, t: tvals[id]})
				// bubble up to keep descending order
				for i := len(displayed) - 1; i > 0 && displayed[i].t > displayed[i-1].t; i-- {
					displayed[i], displayed[i-1] = displayed[i-1], displayed[i]
				}
				inDisp[id] = true
			} else if tvals[id] > displayed[len(displayed)-1].t {
				// replace last
				rem := displayed[len(displayed)-1].id
				inDisp[rem] = false
				displayed[len(displayed)-1] = pair{id: id, t: tvals[id]}
				// bubble up
				for i := len(displayed) - 1; i > 0 && displayed[i].t > displayed[i-1].t; i-- {
					displayed[i], displayed[i-1] = displayed[i-1], displayed[i]
				}
				inDisp[id] = true
			}
		} else {
			if inDisp[id] {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
