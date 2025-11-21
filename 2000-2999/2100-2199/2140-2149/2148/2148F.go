package main

import (
	"bufio"
	"os"
	"strconv"
)

func nextInt(r *bufio.Reader) int {
	sign := 1
	val := 0
	c, err := r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func better(a, b []int, start int) bool {
	la, lb := len(a), len(b)
	i := start
	for i < la && i < lb {
		va, vb := a[i], b[i]
		if va != vb {
			return va < vb
		}
		i++
	}
	if i == la && i == lb {
		return false
	}
	return i == la
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	t := nextInt(reader)
	for ; t > 0; t-- {
		n := nextInt(reader)
		arrays := make([][]int, n)
		maxLen := 0
		for i := 0; i < n; i++ {
			k := nextInt(reader)
			arr := make([]int, k)
			for j := 0; j < k; j++ {
				arr[j] = nextInt(reader)
			}
			arrays[i] = arr
			if k > maxLen {
				maxLen = k
			}
		}

		res := make([]int, maxLen)
		used := make([]bool, n)
		coverage := 0
		for coverage < maxLen {
			best := -1
			for i := 0; i < n; i++ {
				if used[i] || len(arrays[i]) <= coverage {
					continue
				}
				if best == -1 || better(arrays[i], arrays[best], coverage) {
					best = i
				}
			}
			if best == -1 {
				break
			}
			copy(res[coverage:len(arrays[best])], arrays[best][coverage:])
			used[best] = true
			coverage = len(arrays[best])
		}

		for i := 0; i < maxLen; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(strconv.Itoa(res[i]))
		}
		writer.WriteByte('\n')
	}
}
