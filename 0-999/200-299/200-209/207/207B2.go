package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func main() {
	fs := newFastScanner()
	n := fs.nextInt()
	if n == 0 {
		fmt.Println(0)
		return
	}
	orig := make([]int, n)
	for i := 0; i < n; i++ {
		orig[i] = fs.nextInt()
	}

	// Reverse the column so that message constraints become source-based jumps.
	rev := make([]int, n)
	for i := 0; i < n; i++ {
		rev[i] = orig[n-1-i]
	}

	m := 2 * n
	jumps := make([]int, m)
	for i := 0; i < m; i++ {
		jumps[i] = rev[i%n]
	}

	logs := make([]int, m+1)
	for i := 2; i <= m; i++ {
		logs[i] = logs[i/2] + 1
	}
	logM := logs[m] + 1
	st := make([][]int, logM)
	for p := 0; p < logM; p++ {
		size := m - (1 << p) + 1
		if size <= 0 {
			st = st[:p]
			break
		}
		st[p] = make([]int, size)
	}
	logCount := len(st)

	LOG := bits.Len(uint(n)) + 1
	reach := make([][]int, LOG)
	for k := 0; k < LOG; k++ {
		reach[k] = make([]int, m)
	}

	for i := 0; i < m; i++ {
		to := i + jumps[i]
		if to >= m {
			to = m - 1
		}
		reach[0][i] = to
	}

	for k := 0; k < LOG-1; k++ {
		arr := reach[k]
		copy(st[0], arr)
		for p := 1; p < logCount; p++ {
			prev := st[p-1]
			curr := st[p]
			step := 1 << (p - 1)
			limit := len(curr)
			for i := 0; i < limit; i++ {
				curr[i] = max(prev[i], prev[i+step])
			}
		}
		nextArr := reach[k+1]
		for i := 0; i < m; i++ {
			r := arr[i]
			lenRange := r - i + 1
			if lenRange <= 0 {
				lenRange = 1
			}
			p := logs[lenRange]
			best := st[p][i]
			idx := r - (1 << p) + 1
			best = max(best, st[p][idx])
			nextArr[i] = best
		}
	}

	var total int64
	for start := 0; start < n; start++ {
		target := start + n - 1
		pos := start
		if pos >= target {
			continue
		}
		steps := 0
		for k := LOG - 1; k >= 0; k-- {
			nxt := reach[k][pos]
			if nxt < target {
				steps += 1 << k
				pos = nxt
			}
		}
		steps++
		total += int64(steps)
	}

	writer := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(writer, total)
	writer.Flush()
}
