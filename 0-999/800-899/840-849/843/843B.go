package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, start, x int
	if _, err := fmt.Fscan(reader, &n, &start, &x); err != nil {
		return
	}

	rand.Seed(time.Now().UnixNano())

	type Node struct{ val, next int }

	query := func(idx int) Node {
		fmt.Fprintf(writer, "? %d\n", idx)
		writer.Flush()
		var v, nx int
		fmt.Fscan(reader, &v, &nx)
		return Node{v, nx}
	}

	res := query(start)
	if res.val >= x {
		fmt.Fprintf(writer, "! %d\n", res.val)
		writer.Flush()
		return
	}

	best := res
	bestIdx := start

	sample := 1000
	if sample > n {
		sample = n
	}
	for i := 1; i < sample; i++ {
		idx := rand.Intn(n) + 1
		nd := query(idx)
		if nd.val >= x {
			fmt.Fprintf(writer, "! %d\n", nd.val)
			writer.Flush()
			return
		}
		if nd.val > best.val {
			best = nd
			bestIdx = idx
		}
	}

	cur := bestIdx
	for q := sample; q < 1999 && cur != -1; q++ {
		nd := query(cur)
		if nd.val >= x {
			fmt.Fprintf(writer, "! %d\n", nd.val)
			writer.Flush()
			return
		}
		cur = nd.next
	}

	fmt.Fprintf(writer, "! -1\n")
	writer.Flush()
}
