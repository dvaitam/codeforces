package main

import (
	"bufio"
	"fmt"
	"os"
)

func reverseStr(s string) string {
	bs := []byte(s)
	for i, j := 0, len(bs)-1; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}
	return string(bs)
}

func solveCase(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	set2 := make(map[string]bool)
	set3 := make(map[string]bool)

	for _, w := range words {
		if len(w) == 1 {
			fmt.Fprintln(out, "YES")
			return
		}
		rev := reverseStr(w)
		if len(w) == 2 {
			if w[0] == w[1] {
				fmt.Fprintln(out, "YES")
				return
			}
			if set2[rev] || set3[rev] {
				fmt.Fprintln(out, "YES")
				return
			}
			// check length3 words that start with reversed pair
			for c := byte('a'); c <= byte('z'); c++ {
				key := rev + string(c)
				if set3[key] {
					fmt.Fprintln(out, "YES")
					return
				}
			}
			set2[w] = true
		} else if len(w) == 3 {
			if w[0] == w[2] {
				fmt.Fprintln(out, "YES")
				return
			}
			if set3[rev] {
				fmt.Fprintln(out, "YES")
				return
			}
			// check if there was 2-letter word that can form palindrome
			key := rev[:2]
			if set2[key] {
				fmt.Fprintln(out, "YES")
				return
			}
			set3[w] = true
		}
	}
	fmt.Fprintln(out, "NO")
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		solveCase(in, out)
	}
}
