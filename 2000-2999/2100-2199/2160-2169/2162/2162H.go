package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner(reader io.Reader) *FastScanner {
	return &FastScanner{r: bufio.NewReader(reader)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	b, err := fs.r.ReadByte()
	for (b < '0' || b > '9') && b != '-' {
		b, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if b == '-' {
		sign = -1
		b, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for b >= '0' && b <= '9' {
		val = val*10 + int(b-'0')
		b, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	fs := NewFastScanner(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		m := fs.NextInt()

		freq := make([]int, n+2)
		for i := 0; i < n; i++ {
			v := fs.NextInt()
			freq[v]++
		}

		forced := make([]bool, n-1)
		for i := 0; i < m; i++ {
			l := fs.NextInt()
			r := fs.NextInt()
			for pos := l; pos < r; pos++ {
				forced[pos-1] = true
			}
		}

		sizes := make([]int, 0)
		for i := 0; i < n; {
			j := i
			for j+1 < n && forced[j] {
				j++
			}
			sizes = append(sizes, j-i+1)
			i = j + 1
		}

		dp := make([]bool, n+1)
		dp[0] = true
		for _, s := range sizes {
			for sum := n - s; sum >= 0; sum-- {
				if dp[sum] {
					dp[sum+s] = true
				}
			}
		}

		dpPrefix := make([]int, n+2)
		for i := 0; i <= n; i++ {
			dpPrefix[i+1] = dpPrefix[i]
			if dp[i] {
				dpPrefix[i+1]++
			}
		}

		prefFreq := make([]int, n+1)
		for v := 1; v <= n; v++ {
			prefFreq[v] = prefFreq[v-1] + freq[v]
		}

		var sb strings.Builder
		sb.Grow(n + 1)
		for x := 1; x <= n; x++ {
			lowMin := prefFreq[x-1]
			lowMax := prefFreq[x]
			possible := false
			if lowMin <= lowMax {
				if dpPrefix[lowMax+1]-dpPrefix[lowMin] > 0 {
					possible = true
				}
			}
			if possible {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
		sb.WriteByte('\n')
		out.WriteString(sb.String())
	}
}
