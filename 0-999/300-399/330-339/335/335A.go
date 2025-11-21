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

	var s string
	fmt.Fscan(in, &s)
	var n int
	fmt.Fscan(in, &n)

	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}

	if len(freq) > n {
		fmt.Fprintln(out, -1)
		return
	}

	sheets := 1
	base := make(map[rune]int)
	for {
		sum := 0
		for ch, cnt := range freq {
			need := (cnt + sheets - 1) / sheets
			base[ch] = need
			sum += need
		}
		if sum <= n {
			break
		}
		sheets++
	}

	fmt.Fprintln(out, sheets)

	sheet := make([]rune, 0, n)
	for ch := 'a'; ch <= 'z' && len(sheet) < n; ch++ {
		if cnt, ok := base[ch]; ok {
			for i := 0; i < cnt; i++ {
				sheet = append(sheet, ch)
			}
		}
	}

	for len(sheet) < n {
		sheet = append(sheet, 'a')
	}

	fmt.Fprintln(out, string(sheet))
}
