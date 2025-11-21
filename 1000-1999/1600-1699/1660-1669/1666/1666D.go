package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	for ; n > 0; n-- {
		var s, t string
		fmt.Fscan(in, &s, &t)

		var cntS, cntT [26]int
		for i := 0; i < len(s); i++ {
			cntS[s[i]-'A']++
		}
		for i := 0; i < len(t); i++ {
			cntT[t[i]-'A']++
		}

		possible := true
		var removeCnt [26]int
		for i := 0; i < 26; i++ {
			if cntT[i] > cntS[i] {
				possible = false
				break
			}
			removeCnt[i] = cntS[i] - cntT[i]
		}

		if !possible {
			fmt.Fprintln(out, "NO")
			continue
		}

		var removed [26]int
		result := make([]byte, 0, len(s))
		for i := 0; i < len(s); i++ {
			idx := s[i] - 'A'
			if removed[idx] < removeCnt[idx] {
				removed[idx]++
			} else {
				result = append(result, s[i])
			}
		}

		if string(result) == t {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
