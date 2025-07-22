package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	freq := make([]int, 26)
	for _, ch := range s {
		freq[ch-'a']++
	}
	counts := make([]int, 0, 26)
	for _, c := range freq {
		if c > 0 {
			counts = append(counts, c)
		}
	}
	m := len(counts)
	ans := false
	switch {
	case m <= 1:
		ans = false
	case m >= 5:
		ans = false
	case m == 4:
		ans = true
	case m == 3:
		for _, c := range counts {
			if c >= 2 {
				ans = true
				break
			}
		}
	case m == 2:
		if counts[0] >= 2 && counts[1] >= 2 {
			ans = true
		}
	}
	if ans {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
