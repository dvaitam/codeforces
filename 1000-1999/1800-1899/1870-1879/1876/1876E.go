package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Edge represents a directed tree edge with a flag b
type Edge struct{ v, b int }

var (
	n      int
	G, uG  [][]Edge
	f      []int
	ans, R int
	s      []int
	t, col int
	out    *bufio.Writer
)

func d(x, fa int) {
	for _, e := range G[x] {
		v := e.v
		if v != fa {
			d(v, x)
			f[x] += f[v] + 1
		}
	}
	for _, e := range uG[x] {
		v := e.v
		if v != fa {
			d(v, x)
			f[x] += f[v] + (1 - e.b)
		}
	}
}

func d2(x, fa int) {
	if f[x] > ans {
		ans = f[x]
		R = x
	}
	for _, e := range G[x] {
		v, b := e.v, e.b
		if v != fa {
			f[v] = f[x] - 1 + (1 - b)
			d2(v, x)
		}
	}
	for _, e := range uG[x] {
		v, b := e.v, e.b
		if v != fa {
			f[v] = f[x] - (1 - b) + 1
			d2(v, x)
		}
	}
}

func d3(x, fa int) {
	for _, e := range G[x] {
		v := e.v
		if v != fa {
			col++
			t++
			s[t] = col
			fmt.Fprintf(out, "%d %d %d\n", x, v, col)
			d3(v, x)
			t--
		}
	}
	for _, e := range uG[x] {
		v, b := e.v, e.b
		if v != fa {
			if b == 0 {
				col++
				t++
				s[t] = col
				fmt.Fprintf(out, "%d %d %d\n", x, v, col)
				d3(v, x)
				t--
			} else {
				tmp := s[t]
				fmt.Fprintf(out, "%d %d %d\n", v, x, s[t])
				t--
				d3(v, x)
				t++
				s[t] = tmp
			}
		}
	}
}

func readInt(r *bufio.Reader) (int, error) {
	var x, sign int = 0, 1
	// read first byte
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	// skip non-numeric
	for ; (b < '0' || b > '9') && b != '-'; b, err = r.ReadByte() {
		if err != nil {
			return 0, err
		}
	}
	if b == '-' {
		sign = -1
		b, err = r.ReadByte()
		if err != nil {
			return 0, err
		}
	}
	for ; b >= '0' && b <= '9'; b, err = r.ReadByte() {
		if err != nil && err != io.EOF {
			return 0, err
		}
		x = x*10 + int(b-'0')
	}
	return x * sign, nil
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var err error
	n, err = readInt(in)
	if err != nil {
		return
	}
	G = make([][]Edge, n+1)
	uG = make([][]Edge, n+1)
	f = make([]int, n+1)
	s = make([]int, n+1)

	// read edges
	for i := 2; i <= n; i++ {
		x, _ := readInt(in)
		y, _ := readInt(in)
		b, _ := readInt(in)
		G[x] = append(G[x], Edge{v: y, b: b})
		uG[y] = append(uG[y], Edge{v: x, b: b})
	}

	d(1, 0)
	ans = -1
	d2(1, 0)
	fmt.Fprintln(out, ans)
	d3(R, 0)
}
