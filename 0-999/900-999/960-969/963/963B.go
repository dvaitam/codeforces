package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader *bufio.Reader
	writer *bufio.Writer
)

func readInt() int {
	var c byte
	var x int
	var neg bool
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0
		}
		c = b
		if c == '-' {
			neg = true
		}
		if (c >= '0' && c <= '9') || c == '-' {
			break
		}
	}
	if c == '-' {
		neg = true
	} else {
		x = int(c - '0')
	}
	for {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		c = b
		if c < '0' || c > '9' {
			break
		}
		x = x*10 + int(c-'0')
	}
	if neg {
		x = -x
	}
	return x
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	n := readInt()
	N := n + 5
	fir1 := make([]int, N)
	nxt1 := make([]int, 2*N)
	to1 := make([]int, 2*N)
	var cnt1 int

	fir2 := make([]int, N)
	nxt2 := make([]int, 2*N)
	to2 := make([]int, 2*N)
	var cnt2 int
	du := make([]int, N)
	f := make([]int, N)

	add1 := func(x, y int) {
		cnt1++
		nxt1[cnt1] = fir1[x]
		fir1[x] = cnt1
		to1[cnt1] = y
	}
	add2 := func(x, y int) {
		cnt2++
		nxt2[cnt2] = fir2[x]
		fir2[x] = cnt2
		to2[cnt2] = y
		du[y]++
	}

	for i := 1; i <= n; i++ {
		x := readInt()
		if x != 0 {
			add1(i, x)
			add1(x, i)
		}
	}

	var dfs func(x, fa int)
	dfs = func(x, fa int) {
		for e := fir1[x]; e != 0; e = nxt1[e] {
			y := to1[e]
			if y == fa {
				continue
			}
			dfs(y, x)
		}
		if fa != 0 {
			if f[x]&1 == 1 {
				add2(x, fa)
			} else {
				add2(fa, x)
				f[fa]++
			}
		}
	}

	dfs(1, 0)
	if f[1]&1 == 1 {
		fmt.Fprintln(writer, "NO")
		return
	}
	ans := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if du[i] == 0 {
			ans = append(ans, i)
		}
	}
	for i := 0; i < len(ans); i++ {
		u := ans[i]
		for e := fir2[u]; e != 0; e = nxt2[e] {
			v := to2[e]
			du[v]--
			if du[v] == 0 {
				ans = append(ans, v)
			}
		}
	}
	if len(ans) != n {
		fmt.Fprintln(writer, "NO")
		return
	}
	fmt.Fprintln(writer, "YES")
	for _, v := range ans {
		fmt.Fprintln(writer, v)
	}
}
