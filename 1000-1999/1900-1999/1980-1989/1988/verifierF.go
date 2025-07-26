package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 998244353

func upd(x *int64, y int64) { *x = (*x + y) % mod }

func solve(n int, a, b, c []int) []int64 {
	C := make([][]int, n+1)
	for i := range C {
		C[i] = make([]int, n+1)
	}
	C[0][0] = 1
	for i := 1; i <= n; i++ {
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			C[i][j] = (C[i-1][j] + C[i-1][j-1]) % mod
		}
	}
	f := make([][][]int64, 2)
	for i := range f {
		f[i] = make([][]int64, n+2)
		for j := range f[i] {
			f[i][j] = make([]int64, n+2)
		}
	}
	pre := make([][]int64, n+1)
	suf := make([][]int64, n+1)
	g := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		pre[i] = make([]int64, n+1)
		suf[i] = make([]int64, n+1)
		g[i] = make([]int64, n+1)
	}
	f[1][1][0] = 1
	pre[0][0] = int64(a[0])
	suf[0][0] = int64(b[0])
	for i := 1; i <= n; i++ {
		nw := i & 1
		nx := nw ^ 1
		for j := range f[nx] {
			for k := range f[nx][j] {
				f[nx][j][k] = 0
			}
		}
		for j := 1; j <= i; j++ {
			for k := 0; k <= i-1; k++ {
				if f[nw][j][k] != 0 {
					val := f[nw][j][k]
					upd(&pre[i][k], val*int64(a[j]))
					upd(&suf[i][i-1-k], val*int64(b[j]))
					upd(&f[nx][j+1][k+1], val)
					upd(&f[nx][j][k], val*int64(k+1))
					upd(&f[nx][j][k+1], val*int64(i-k-1))
				}
			}
		}
	}
	for i := 0; i <= n; i++ {
		for x := 0; x <= i; x++ {
			for y := 0; y <= n-x-1; y++ {
				t := 0
				if i > 0 {
					t = 1
				}
				upd(&g[i][y], pre[i][x]*int64(c[x+y+t]))
			}
		}
	}
	ans := make([]int64, n)
	for i := 1; i <= n; i++ {
		var cur int64
		for p := 1; p <= i; p++ {
			for y := 0; y <= i-p; y++ {
				upd(&cur, g[p-1][y]*suf[i-p][y]%mod*int64(C[i-1][p-1])%mod)
			}
		}
		ans[i-1] = cur
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	for t := 0; t < 100; t++ {
		n := rand.Intn(6) + 1
		a := make([]int, n)
		b := make([]int, n)
		c := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rand.Intn(5) + 1
		}
		for i := 0; i < n; i++ {
			b[i] = rand.Intn(5) + 1
		}
		for i := 0; i < n; i++ {
			c[i] = rand.Intn(5) + 1
		}
		input := &bytes.Buffer{}
		fmt.Fprintf(input, "%d\n", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(input, "%d ", a[i])
		}
		fmt.Fprintln(input)
		for i := 0; i < n; i++ {
			fmt.Fprintf(input, "%d ", b[i])
		}
		fmt.Fprintln(input)
		for i := 0; i < n; i++ {
			fmt.Fprintf(input, "%d ", c[i])
		}
		fmt.Fprintln(input)
		expected := solve(n, a, b, c)
		var expBuf bytes.Buffer
		for i, v := range expected {
			if i+1 == len(expected) {
				fmt.Fprintf(&expBuf, "%d ", v%mod)
			} else {
				fmt.Fprintf(&expBuf, "%d ", v%mod)
			}
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("error running binary:", err)
			os.Exit(1)
		}
		if strings.TrimSpace(out.String()) != strings.TrimSpace(expBuf.String()) {
			fmt.Println("wrong answer on test", t+1)
			fmt.Println("expected:", expBuf.String())
			fmt.Println("got:", out.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
