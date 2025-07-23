package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	freq := make([]int, 26)
	for i := 0; i < n; i++ {
		freq[s[i]-'a']++
	}
	remove := make([]int, 26)
	for c := 0; c < 26 && k > 0; c++ {
		if freq[c] <= k {
			remove[c] = freq[c]
			k -= freq[c]
		} else {
			remove[c] = k
			k = 0
		}
	}
	var res []byte
	res = make([]byte, 0, n)
	for i := 0; i < n; i++ {
		ch := s[i] - 'a'
		if remove[ch] > 0 {
			remove[ch]--
		} else {
			res = append(res, s[i])
		}
	}
	writer.Write(res)
}
