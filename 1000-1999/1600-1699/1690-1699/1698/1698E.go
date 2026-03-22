package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

const MOD int64 = 998244353

type FastScanner struct {
	data []byte
	idx  int
}

func (fs *FastScanner) NextInt() int {
	for fs.idx < len(fs.data) {
		c := fs.data[fs.idx]
		if c != ' ' && c != '\n' && c != '\r' && c != '\t' {
			break
		}
		fs.idx++
	}
	sign := 1
	if fs.data[fs.idx] == '-' {
		sign = -1
		fs.idx++
	}
	val := 0
	for fs.idx < len(fs.data) {
		c := fs.data[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		val = val*10 + int(c-'0')
		fs.idx++
	}
	return sign * val
}

func main() {
	data, _ := io.ReadAll(os.Stdin)
	fs := FastScanner{data: data}
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		s := fs.NextInt()

		a := make([]int, n+1)
		posA := make([]int, n+1)
		for i := 1; i <= n; i++ {
			a[i] = fs.NextInt()
			posA[a[i]] = i
		}

		b := make([]int, n+1)
		for i := 1; i <= n; i++ {
			b[i] = fs.NextInt()
		}

		used := make([]bool, n+1)
		missPos := make([]int, 0, n)
		ok := true

		for i := 1; i <= n; i++ {
			x := b[posA[i]]
			if x == -1 {
				missPos = append(missPos, i)
			} else {
				if used[x] {
					ok = false
				}
				used[x] = true
				if i > x+s {
					ok = false
				}
			}
		}

		if !ok {
			out.WriteString("0\n")
			continue
		}

		suf := make([]int, n+2)
		for v := n; v >= 1; v-- {
			suf[v] = suf[v+1]
			if !used[v] {
				suf[v]++
			}
		}

		var ans int64 = 1
		assigned := 0
		for i := len(missPos) - 1; i >= 0; i-- {
			pos := missPos[i]
			thr := pos - s
			if thr < 1 {
				thr = 1
			}
			choices := suf[thr] - assigned
			if choices <= 0 {
				ans = 0
				break
			}
			ans = ans * int64(choices) % MOD
			assigned++
		}

		out.WriteString(strconv.FormatInt(ans, 10))
		out.WriteByte('\n')
	}
}