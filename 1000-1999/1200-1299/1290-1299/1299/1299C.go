package main

import (
	"bufio"
	"fmt"
	"os"
)

type segment struct {
	sum int64
	cnt int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	segs := make([]segment, 0, n)
	for i := 0; i < n; i++ {
		var v int64
		fmt.Fscan(in, &v)
		segs = append(segs, segment{sum: v, cnt: 1})
		for len(segs) >= 2 {
			m := len(segs)
			a := segs[m-2]
			b := segs[m-1]
			if a.sum*int64(b.cnt) > b.sum*int64(a.cnt) {
				segs[m-2].sum += b.sum
				segs[m-2].cnt += b.cnt
				segs = segs[:m-1]
			} else {
				break
			}
		}
	}
	for _, s := range segs {
		avg := float64(s.sum) / float64(s.cnt)
		for i := 0; i < s.cnt; i++ {
			fmt.Fprintf(out, "%.10f\n", avg)
		}
	}
}
