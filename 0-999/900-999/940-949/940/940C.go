package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	var s string
	fmt.Fscan(in, &s)

	letterSet := make(map[byte]bool)
	for i := 0; i < len(s); i++ {
		letterSet[s[i]] = true
	}
	letters := make([]byte, 0, len(letterSet))
	for ch := range letterSet {
		letters = append(letters, ch)
	}
	sort.Slice(letters, func(i, j int) bool { return letters[i] < letters[j] })

	minChar := letters[0]

	if k > n {
		res := make([]byte, 0, k)
		res = append(res, s...)
		for i := 0; i < k-n; i++ {
			res = append(res, minChar)
		}
		fmt.Println(string(res))
		return
	}

	t := []byte(s[:k])
	for i := k - 1; i >= 0; i-- {
		curr := t[i]
		idx := sort.Search(len(letters), func(j int) bool { return letters[j] > curr })
		if idx < len(letters) {
			t[i] = letters[idx]
			for j := i + 1; j < k; j++ {
				t[j] = minChar
			}
			fmt.Println(string(t))
			return
		}
	}
}
