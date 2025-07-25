package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

type Query struct {
	s, d, k int
}

type Precalc struct {
	pre  [][]int64
	wpre [][]int64
}

func buildPrecalc(a []int64, n, d int) Precalc {
	pre := make([][]int64, d)
	wpre := make([][]int64, d)
	for r := 0; r < d; r++ {
		length := (n-r-1)/d + 1
		pre[r] = make([]int64, length+1)
		wpre[r] = make([]int64, length+1)
		var sum int64
		var wsum int64
		for j := 0; j < length; j++ {
			idx := r + 1 + j*d
			val := a[idx]
			sum += val
			pre[r][j+1] = sum
			wsum += val * int64(j+1)
			wpre[r][j+1] = wsum
		}
	}
	return Precalc{pre: pre, wpre: wpre}
}

func querySmall(pc Precalc, s, d, k int) int64 {
	r := (s - 1) % d
	idx := (s - 1 - r) / d
	pre := pc.pre[r]
	wpre := pc.wpre[r]
	res := wpre[idx+k] - wpre[idx] - int64(idx)*(pre[idx+k]-pre[idx])
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		a := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		queries := make([]Query, q)
		smallSet := make(map[int]struct{})
		threshold := int(math.Sqrt(float64(n))) + 1
		for i := 0; i < q; i++ {
			fmt.Fscan(reader, &queries[i].s, &queries[i].d, &queries[i].k)
			if queries[i].d <= threshold {
				smallSet[queries[i].d] = struct{}{}
			}
		}
		precalc := make(map[int]Precalc)
		for d := range smallSet {
			precalc[d] = buildPrecalc(a, n, d)
		}
		answers := make([]int64, q)
		for i, qu := range queries {
			if qu.d <= threshold {
				pc := precalc[qu.d]
				answers[i] = querySmall(pc, qu.s, qu.d, qu.k)
			} else {
				var sum int64
				pos := qu.s
				for j := 0; j < qu.k; j++ {
					sum += int64(j+1) * a[pos]
					pos += qu.d
				}
				answers[i] = sum
			}
		}
		for i, ans := range answers {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, ans)
		}
		fmt.Fprintln(writer)
	}
}
