package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
)

type FastScanner struct {
	data []byte
	idx  int
	n    int
}

func NewFastScanner() *FastScanner {
	data, _ := io.ReadAll(os.Stdin)
	return &FastScanner{data: data, n: len(data)}
}

func (fs *FastScanner) skip() {
	for fs.idx < fs.n && fs.data[fs.idx] <= ' ' {
		fs.idx++
	}
}

func (fs *FastScanner) NextInt() int {
	fs.skip()
	sign := 1
	if fs.data[fs.idx] == '-' {
		sign = -1
		fs.idx++
	}
	val := 0
	for fs.idx < fs.n && fs.data[fs.idx] > ' ' {
		val = val*10 + int(fs.data[fs.idx]-'0')
		fs.idx++
	}
	return val * sign
}

func (fs *FastScanner) NextBytes() []byte {
	fs.skip()
	start := fs.idx
	for fs.idx < fs.n && fs.data[fs.idx] > ' ' {
		fs.idx++
	}
	return fs.data[start:fs.idx]
}

func writeInt(out *bytes.Buffer, x int) {
	out.WriteString(strconv.Itoa(x))
}

func main() {
	fs := NewFastScanner()
	t := fs.NextInt()

	var out bytes.Buffer

	for ; t > 0; t-- {
		n := fs.NextInt()
		m := fs.NextInt()

		g := make([][]byte, n)
		for i := 0; i < n; i++ {
			g[i] = fs.NextBytes()
		}

		const inf = int(1e9)
		var minR, maxR, minC, maxC [26]int
		for i := 0; i < 26; i++ {
			minR[i], minC[i] = inf, inf
			maxR[i], maxC[i] = -1, -1
		}

		mx := -1
		for r := 0; r < n; r++ {
			row := g[r]
			for c := 0; c < m; c++ {
				ch := row[c]
				if ch == '.' {
					continue
				}
				id := int(ch - 'a')
				if id > mx {
					mx = id
				}
				if r < minR[id] {
					minR[id] = r
				}
				if r > maxR[id] {
					maxR[id] = r
				}
				if c < minC[id] {
					minC[id] = c
				}
				if c > maxC[id] {
					maxC[id] = c
				}
			}
		}

		if mx == -1 {
			out.WriteString("YES\n0\n")
			continue
		}

		coords := make([][4]int, mx+1)
		ok := true

		for id := 0; id <= mx && ok; id++ {
			if maxR[id] == -1 {
				continue
			}

			if minR[id] != maxR[id] && minC[id] != maxC[id] {
				ok = false
				break
			}

			if minR[id] == maxR[id] {
				r := minR[id]
				c1, c2 := minC[id], maxC[id]
				for c := c1; c <= c2; c++ {
					ch := g[r][c]
					if ch == '.' || int(ch-'a') < id {
						ok = false
						break
					}
				}
				coords[id] = [4]int{r, c1, r, c2}
			} else {
				c := minC[id]
				r1, r2 := minR[id], maxR[id]
				for r := r1; r <= r2; r++ {
					ch := g[r][c]
					if ch == '.' || int(ch-'a') < id {
						ok = false
						break
					}
				}
				coords[id] = [4]int{r1, c, r2, c}
			}
		}

		if !ok {
			out.WriteString("NO\n")
			continue
		}

		last := coords[mx]
		for id := mx; id >= 0; id-- {
			if maxR[id] == -1 {
				coords[id] = last
			} else {
				last = coords[id]
			}
		}

		out.WriteString("YES\n")
		writeInt(&out, mx+1)
		out.WriteByte('\n')
		for i := 0; i <= mx; i++ {
			c := coords[i]
			writeInt(&out, c[0]+1)
			out.WriteByte(' ')
			writeInt(&out, c[1]+1)
			out.WriteByte(' ')
			writeInt(&out, c[2]+1)
			out.WriteByte(' ')
			writeInt(&out, c[3]+1)
			out.WriteByte('\n')
		}
	}

	w := bufio.NewWriterSize(os.Stdout, 1<<20)
	w.Write(out.Bytes())
	w.Flush()
}