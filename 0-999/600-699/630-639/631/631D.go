package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type segment struct {
	len  int64
	char byte
}

func readSegments(reader *bufio.Reader, count int) []segment {
	segs := make([]segment, 0, count)
	for i := 0; i < count; i++ {
		var token string
		fmt.Fscan(reader, &token)
		parts := strings.Split(token, "-")
		l, _ := strconv.ParseInt(parts[0], 10, 64)
		c := parts[1][0]
		if len(segs) > 0 && segs[len(segs)-1].char == c {
			segs[len(segs)-1].len += l
		} else {
			segs = append(segs, segment{len: l, char: c})
		}
	}
	return segs
}

func equal(a, b segment) bool {
	return a.len == b.len && a.char == b.char
}

func prefixFunction(p []segment) []int {
	pi := make([]int, len(p))
	for i := 1; i < len(p); i++ {
		j := pi[i-1]
		for j > 0 && !equal(p[i], p[j]) {
			j = pi[j-1]
		}
		if equal(p[i], p[j]) {
			j++
		}
		pi[i] = j
	}
	return pi
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	tSegs := readSegments(reader, n)
	sSegs := readSegments(reader, m)

	m = len(sSegs)
	n = len(tSegs)

	if m == 1 {
		ans := int64(0)
		pat := sSegs[0]
		for _, seg := range tSegs {
			if seg.char == pat.char && seg.len >= pat.len {
				ans += seg.len - pat.len + 1
			}
		}
		fmt.Fprintln(writer, ans)
		return
	}

	if m == 2 {
		ans := 0
		for i := 0; i < n-1; i++ {
			if tSegs[i].char == sSegs[0].char && tSegs[i+1].char == sSegs[1].char &&
				tSegs[i].len >= sSegs[0].len && tSegs[i+1].len >= sSegs[1].len {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
		return
	}

	// m >= 3
	midPattern := sSegs[1 : m-1]
	pi := prefixFunction(midPattern)
	ans := 0
	for i, j := 0, 0; i < n; i++ {
		for j > 0 && !equal(tSegs[i], midPattern[j]) {
			j = pi[j-1]
		}
		if equal(tSegs[i], midPattern[j]) {
			j++
			if j == len(midPattern) {
				// candidate starting index in tSegs for middle match
				start := i - len(midPattern) + 1
				if start-1 >= 0 && start+len(midPattern) < n {
					left := tSegs[start-1]
					right := tSegs[start+len(midPattern)]
					if left.char == sSegs[0].char && left.len >= sSegs[0].len &&
						right.char == sSegs[m-1].char && right.len >= sSegs[m-1].len {
						ans++
					}
				}
				j = pi[j-1]
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
