package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	MAXQ = 400000 + 5
	LOG  = 20
	INF  = int64(1 << 60)
)

var (
	w   [MAXQ]int64
	up  [MAXQ][LOG]int
	sum [MAXQ][LOG]int64
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var Q int
	if _, err := fmt.Fscan(reader, &Q); err != nil {
		return
	}

	// initialize virtual node 0
	w[0] = INF
	for i := 0; i < LOG; i++ {
		sum[0][i] = INF
	}

	cnt := 1
	w[1] = 0
	for i := 0; i < LOG; i++ {
		up[1][i] = 0
		sum[1][i] = INF
	}

	last := 0
	for ; Q > 0; Q-- {
		var typ int
		fmt.Fscan(reader, &typ)
		if typ == 1 {
			var pIn, qIn int64
			fmt.Fscan(reader, &pIn, &qIn)
			r := int(pIn ^ int64(last))
			weight := qIn ^ int64(last)
			cnt++
			w[cnt] = weight
			cur := r
			if w[cur] < weight {
				for i := LOG - 1; i >= 0; i-- {
					anc := up[cur][i]
					if anc != 0 && w[anc] < weight {
						cur = anc
					}
				}
				cur = up[cur][0]
			}
			up[cnt][0] = cur
			if cur == 0 {
				sum[cnt][0] = INF
			} else {
				sum[cnt][0] = w[cur]
			}
			for i := 1; i < LOG; i++ {
				up[cnt][i] = up[up[cnt][i-1]][i-1]
				s := sum[cnt][i-1] + sum[up[cnt][i-1]][i-1]
				if s > INF {
					s = INF
				}
				sum[cnt][i] = s
			}
		} else if typ == 2 {
			var pIn, qIn int64
			fmt.Fscan(reader, &pIn, &qIn)
			r := int(pIn ^ int64(last))
			X := qIn ^ int64(last)
			if w[r] > X {
				last = 0
				fmt.Fprintln(writer, 0)
				continue
			}
			X -= w[r]
			ans := 1
			cur := r
			for i := LOG - 1; i >= 0; i-- {
				if up[cur][i] != 0 && sum[cur][i] <= X {
					X -= sum[cur][i]
					cur = up[cur][i]
					ans += 1 << uint(i)
				}
			}
			last = ans
			fmt.Fprintln(writer, ans)
		}
	}
}
