package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// transform applies one operation to x by appending bit b ('0' or '1')
// to its binary representation, reversing the resulting string,
// removing leading zeros and converting back to an integer.
func transform(x int64, b byte) int64 {
	s := strconv.FormatInt(x, 2)
	s += string(b)
	bs := []byte(s)
	for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}
	i := 0
	for i < len(bs) && bs[i] == '0' {
		i++
	}
	if i == len(bs) {
		return 0
	}
	val, _ := strconv.ParseInt(string(bs[i:]), 2, 64)
	return val
}

func main() {
	br := bufio.NewReader(os.Stdin)
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	var x, y int64
	fmt.Fscan(br, &x, &y)

	if x == y {
		fmt.Fprintln(bw, "YES")
		return
	}

	limit := int64(2)
	if x > y {
		limit *= x
	} else {
		limit *= y
	}

	queue := []int64{x}
	visited := map[int64]bool{x: true}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur == y {
			fmt.Fprintln(bw, "YES")
			return
		}
		if cur > limit {
			continue
		}
		for _, bit := range []byte{'0', '1'} {
			nxt := transform(cur, bit)
			if nxt <= limit && !visited[nxt] {
				visited[nxt] = true
				queue = append(queue, nxt)
			}
		}
	}

	fmt.Fprintln(bw, "NO")
}
