package main

import (
	"bufio"
	"os"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	b, err := fs.r.ReadByte()
	for err == nil && (b <= ' ' || b == '\n') {
		b, err = fs.r.ReadByte()
	}
	if err != nil {
		return 0
	}
	if b == '-' {
		sign = -1
		b, err = fs.r.ReadByte()
	}
	for err == nil && b >= '0' && b <= '9' {
		val = val*10 + int(b-'0')
		b, err = fs.r.ReadByte()
	}
	return sign * val
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.nextInt()
	for ; t > 0; t-- {
		n := fs.nextInt()
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(fs.nextInt())
		}

		rows := make([][]int64, n)
		for i := 0; i < n; i++ {
			row := make([]int64, n)
			row[i] = 1
			rows[i] = row
		}

		var best int64
		first := true
		for len(rows) > 0 {
			sumVec := make([]int64, n)
			for _, row := range rows {
				for j := 0; j < n; j++ {
					sumVec[j] += row[j]
				}
			}
			var dot int64
			for j := 0; j < n; j++ {
				dot += sumVec[j] * a[j]
			}
			val := abs64(dot)
			if first || val > best {
				best = val
				first = false
			}
			if len(rows) == 1 {
				break
			}
			nextRows := make([][]int64, len(rows)-1)
			for i := 0; i < len(rows)-1; i++ {
				row := make([]int64, n)
				for j := 0; j < n; j++ {
					row[j] = rows[i+1][j] - rows[i][j]
				}
				nextRows[i] = row
			}
			rows = nextRows
		}

		out.WriteString(int64ToString(best))
		out.WriteByte('\n')
	}
}

func int64ToString(x int64) string {
	if x == 0 {
		return "0"
	}
	var buf [32]byte
	idx := len(buf)
	neg := x < 0
	if neg {
		x = -x
	}
	for x > 0 {
		idx--
		buf[idx] = byte('0' + x%10)
		x /= 10
	}
	if neg {
		idx--
		buf[idx] = '-'
	}
	return string(buf[idx:])
}
