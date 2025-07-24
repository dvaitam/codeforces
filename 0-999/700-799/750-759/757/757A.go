package main

import "fmt"

func main() {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return
	}
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}
	need := map[rune]int{'B': 1, 'u': 2, 'l': 1, 'b': 1, 'a': 2, 's': 1, 'r': 1}
	ans := int(^uint(0) >> 1)
	for ch, cnt := range need {
		if v, ok := freq[ch]; ok {
			if v/cnt < ans {
				ans = v / cnt
			}
		} else {
			ans = 0
			break
		}
	}
	fmt.Print(ans)
}
