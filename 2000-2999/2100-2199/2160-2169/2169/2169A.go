package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxB int64 = 2_000_000_000

type fastScanner struct {
	r *bufio.Reader
}

func newFastScanner() *fastScanner {
	return &fastScanner{
		r: bufio.NewReader(os.Stdin),
	}
}

func (fs *fastScanner) nextInt64() int64 {
	var sign int64 = 1
	var val int64
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int64(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func main() {
	sc := newFastScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := int(sc.nextInt64())
	for ; t > 0; t-- {
		n := int(sc.nextInt64())
		a := sc.nextInt64()

		events := make(map[int64]int, 2*n+4)
		for i := 0; i < n; i++ {
			v := sc.nextInt64()
			dist := v - a
			if dist < 0 {
				dist = -dist
			}
			if dist == 0 {
				continue
			}
			left := v - dist + 1
			right := v + dist - 1
			if left < 0 {
				left = 0
			}
			if right > maxB {
				right = maxB
			}
			if left > right {
				continue
			}
			events[left]++
			events[right+1]--
		}

		if len(events) == 0 {
			fmt.Fprintln(out, 0)
			continue
		}

		keys := make([]int64, 0, len(events))
		for k := range events {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

		cur := 0
		bestCount := -1
		bestB := int64(0)

		for _, key := range keys {
			cur += events[key]
			if key > maxB {
				continue
			}
			if cur > bestCount {
				bestCount = cur
				bestB = key
			}
		}

		if bestCount < 0 || bestB < 0 {
			bestB = 0
		}

		fmt.Fprintln(out, bestB)
	}
}
