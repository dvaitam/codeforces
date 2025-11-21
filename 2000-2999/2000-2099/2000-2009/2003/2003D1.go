package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() (int, error) {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0, err
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0, err
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val, nil
}

func main() {
	fs := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t, _ := fs.nextInt()
	for ; t > 0; t-- {
		n, _ := fs.nextInt()
		mVal, _ := fs.nextInt()
		m := int64(mVal)

		maxC := 0
		for i := 0; i < n; i++ {
			l, _ := fs.nextInt()
			arr := make([]int, l)
			for j := 0; j < l; j++ {
				arr[j], _ = fs.nextInt()
			}
			sort.Ints(arr)
			uniq := arr[:0]
			for _, v := range arr {
				if len(uniq) == 0 || v != uniq[len(uniq)-1] {
					uniq = append(uniq, v)
				}
			}

			b := 0
			idx := 0
			for idx < len(uniq) && uniq[idx] == b {
				b++
				idx++
			}

			c := b + 1
			for idx < len(uniq) && uniq[idx] == c {
				c++
				idx++
			}

			if c > maxC {
				maxC = c
			}
		}

		M := int64(maxC)
		var ans int64
		if m <= M {
			ans = (m + 1) * M
		} else {
			part1 := (M + 1) * M
			part2 := (m + M + 1) * (m - M) / 2
			ans = part1 + part2
		}

		fmt.Fprintln(out, ans)
	}
}
