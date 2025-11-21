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
	c, _ := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, _ = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, _ = fs.r.ReadByte()
	}
	return sign * val
}

const negInf = int64(-1 << 60)

func main() {
	in := NewFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.NextInt()
	for ; t > 0; t-- {
		n := in.NextInt()
		m := in.NextInt()
		_ = in.NextInt() // k, unused directly

		freq := make([]map[int]int, n+2)
		blanks := make([]int, n+2)
		for i := 1; i <= n; i++ {
			for j := 0; j < m; j++ {
				val := in.NextInt()
				if val == -1 {
					blanks[i]++
				} else {
					if freq[i] == nil {
						freq[i] = make(map[int]int)
					}
					freq[i][val]++
				}
			}
		}

		constSum := int64(0)
		for i := 1; i < n; i++ {
			fi := freq[i]
			fj := freq[i+1]
			if fi == nil || fj == nil || len(fi) == 0 || len(fj) == 0 {
				continue
			}
			if len(fi) > len(fj) {
				fi, fj = fj, fi
			}
			for color, cnt := range fi {
				if cnt2, ok := fj[color]; ok {
					constSum += int64(cnt) * int64(cnt2)
				}
			}
		}

	dpPrev := make(map[int]int64)
	dpPrev[0] = 0

	for i := 1; i <= n; i++ {
		best1 := negInf
		best2 := negInf
		bestColor := 0
		for color, val := range dpPrev {
			if val > best1 {
				best2 = best1
				best1 = val
				bestColor = color
			} else if val > best2 {
				best2 = val
			}
		}

		candMap := make(map[int]struct{})
		candList := make([]int, 0)
		add := func(color int) {
			if _, ok := candMap[color]; ok {
				return
			}
			candMap[color] = struct{}{}
			candList = append(candList, color)
		}
		add(0)
		if i > 1 && freq[i-1] != nil {
			for color := range freq[i-1] {
				add(color)
			}
		}
		if i < n && freq[i+1] != nil {
			for color := range freq[i+1] {
				add(color)
			}
		}

		dpCurr := make(map[int]int64, len(candList))
		JPrev := int64(blanks[i-1]) * int64(blanks[i])
		for _, color := range candList {
			prevVal, ok := dpPrev[color]
			if !ok {
				prevVal = dpPrev[0]
			}
			bestDifferent := best1
			if ok && color == bestColor {
				bestDifferent = best2
			}

			same := prevVal + JPrev
			best := bestDifferent
			if same > best {
				best = same
			}

			left := 0
			if freq[i-1] != nil {
				if v, exists := freq[i-1][color]; exists {
					left = v
				}
			}
			right := 0
			if freq[i+1] != nil {
				if v, exists := freq[i+1][color]; exists {
					right = v
				}
			}
			gain := int64(blanks[i]) * int64(left+right)
			dpCurr[color] = gain + best
		}

		dpPrev = dpCurr
	}

	bestFinal := negInf
	for _, val := range dpPrev {
		if val > bestFinal {
			bestFinal = val
		}
	}
	fmt.Fprintln(out, constSum+bestFinal)
}
