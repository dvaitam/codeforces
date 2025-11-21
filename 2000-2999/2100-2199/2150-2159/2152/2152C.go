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

func (fs *FastScanner) NextInt() int {
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
	return val * sign
}

func main() {
	fs := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := fs.NextInt()
	for ; t > 0; t-- {
		n := fs.NextInt()
		q := fs.NextInt()

		arr := make([]int, n+1)
		for i := 1; i <= n; i++ {
			arr[i] = fs.NextInt()
		}

		zeros := make([]int, n+1)
		for i := 1; i <= n; i++ {
			zeros[i] = zeros[i-1]
			if arr[i] == 0 {
				zeros[i]++
			}
		}

		equal := make([]int, n+1)
		for i := 1; i < n; i++ {
			equal[i] = equal[i-1]
			if arr[i] == arr[i+1] {
				equal[i]++
			}
		}
		equal[n] = equal[n-1]

		for ; q > 0; q-- {
			l := fs.NextInt()
			r := fs.NextInt()

			length := r - l + 1
			countZero := zeros[r] - zeros[l-1]
			countOne := length - countZero

			if countZero%3 != 0 || countOne%3 != 0 {
				fmt.Fprintln(out, -1)
				continue
			}

			ans := length / 3
			if equal[r-1]-equal[l-1] == 0 {
				ans++
			}
			fmt.Fprintln(out, ans)
		}
	}
}
