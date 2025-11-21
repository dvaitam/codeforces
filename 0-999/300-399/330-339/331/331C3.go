package main

import (
	"bufio"
	"fmt"
	"os"
)

type answer struct {
	step int64
	rem  int64
}

var (
	tp   [20]int64
	memo [10][20][11]answer
)

func initMemo() {
	for h := 0; h < 10; h++ {
		for n := 0; n < 20; n++ {
			for idx := 0; idx < 11; idx++ {
				memo[h][n][idx] = answer{step: -1, rem: 0}
			}
		}
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func query(h, n int, x int64) answer {
	idx := 10
	if n >= 1 && x/10 == tp[n-1]-1 {
		idx = int(x % 10)
	}
	key := &memo[h][n][idx]
	if key.step != -1 {
		return *key
	}
	if n == 1 {
		if x >= int64(h) {
			*key = answer{step: 1, rem: 0}
		} else {
			*key = answer{step: 0, rem: x}
		}
		return *key
	}

	cur := x % tp[n-1]
	t := int(x / tp[n-1])
	var cnt int64
	for t >= 0 {
		tmp := query(maxInt(h, t), n-1, cur)
		cnt += tmp.step
		cur = tmp.rem
		if t > 0 {
			cur = cur + tp[n-1] - int64(maxInt(h, t))
			cnt++
		}
		t--
	}
	*key = answer{step: cnt, rem: cur}
	return *key
}

func main() {
	tp[0] = 1
	for i := 1; i <= 18; i++ {
		tp[i] = tp[i-1] * 10
	}
	tp[19] = 9000000000000000000
	initMemo()

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	if n == 0 {
		fmt.Fprintln(out, 0)
		return
	}
	p := 0
	for p < len(tp) && tp[p] <= n {
		p++
	}
	ans := query(0, p, n)
	fmt.Fprintln(out, ans.step)
}
