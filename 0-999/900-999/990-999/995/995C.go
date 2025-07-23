package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const limit = 1500000
const limitSq = int64(limit) * int64(limit)

type Vec struct {
	x, y int64
	idx  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	vs := make([]Vec, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &vs[i].x, &vs[i].y)
		vs[i].idx = i
	}

	rand.Seed(time.Now().UnixNano())
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	signs := make([]int, n)

	for {
		rand.Shuffle(n, func(i, j int) { order[i], order[j] = order[j], order[i] })
		sx, sy := int64(0), int64(0)
		for _, id := range order {
			v := vs[id]
			sx1, sy1 := sx+v.x, sy+v.y
			sx2, sy2 := sx-v.x, sy-v.y
			if sx1*sx1+sy1*sy1 <= sx2*sx2+sy2*sy2 {
				signs[v.idx] = 1
				sx, sy = sx1, sy1
			} else {
				signs[v.idx] = -1
				sx, sy = sx2, sy2
			}
		}
		if sx*sx+sy*sy <= limitSq {
			writer := bufio.NewWriter(os.Stdout)
			defer writer.Flush()
			for i := 0; i < n; i++ {
				if i > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprintf(writer, "%d", signs[i])
			}
			writer.WriteByte('\n')
			return
		}
	}
}
