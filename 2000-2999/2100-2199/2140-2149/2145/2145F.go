package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *FastScanner) NextInt64() int64 {
	sign := int64(1)
	val := int64(0)
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return val * sign
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func lcm(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	return a / gcd(a, b) * b
}

type stateInfo struct {
	turn int64
	pos  int64
}

func solveCase(n int, m int64, a, b []int64) int64 {
	l := int64(1)
	for i := 0; i < n; i++ {
		l = lcm(l, a[i])
	}

	pos := int64(0)
	turn := int64(1)
	startIdx := -1
	// Remember states (pos modulo L, active trap) to fast-forward repeating patterns.
	visited := make(map[int64]stateInfo)

	for {
		if startIdx >= 0 {
			idx := startIdx
			if pos%a[idx] == b[idx] {
				if pos == 0 {
					return -1
				}
				pos--
			}
			key := int64(idx)*l + pos%l
			if info, ok := visited[key]; ok {
				if info.turn < turn {
					cycleLen := turn - info.turn
					cycleAdvance := pos - info.pos
					if cycleAdvance == 0 {
						return -1
					}
					remaining := m - pos
					if remaining > 0 {
						k := (remaining - 1) / cycleAdvance
						if k > 0 {
							// Skip k whole cycles since behavior repeats.
							pos += cycleAdvance * k
							turn += cycleLen * k
							visited[key] = stateInfo{turn: turn, pos: pos}
							continue
						}
					}
				}
			} else {
				visited[key] = stateInfo{turn: turn, pos: pos}
			}
		}

		pos++
		var endIdx int
		if startIdx == -1 {
			endIdx = 0
		} else {
			endIdx = (startIdx + 1) % n
		}
		if pos == m && m%a[endIdx] != b[endIdx] {
			return turn
		}

		turn++
		if startIdx == -1 {
			startIdx = 0
		} else {
			startIdx = (startIdx + 1) % n
		}
	}
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(fs.NextInt64())
	for ; t > 0; t-- {
		n := int(fs.NextInt64())
		m := fs.NextInt64()
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = fs.NextInt64()
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			b[i] = fs.NextInt64()
		}
		fmt.Fprintln(out, solveCase(n, m, a, b))
	}
}
