package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func buildSong(s0, t string, k int) string {
	s := s0
	for i := 0; i < k && i < len(t); i++ {
		s = s + string(t[i]) + s
		if len(s) > 2000 {
			break
		}
	}
	return s
}

func countOccurrences(s, w string) int64 {
	if len(w) == 0 || len(w) > len(s) {
		return 0
	}
	count := int64(0)
	for i := 0; i+len(w) <= len(s); i++ {
		if s[i:i+len(w)] == w {
			count++
		}
	}
	return count % mod
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	var s0, t string
	fmt.Fscan(reader, &s0)
	fmt.Fscan(reader, &t)

	for i := 0; i < q; i++ {
		var k int
		var w string
		fmt.Fscan(reader, &k, &w)
		song := buildSong(s0, t, k)
		cnt := countOccurrences(song, w)
		fmt.Fprintln(writer, cnt%mod)
	}
}
