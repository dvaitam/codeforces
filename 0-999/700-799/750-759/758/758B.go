package main

import "fmt"

func main() {
	var s string
	if _, err := fmt.Scan(&s); err != nil {
		return
	}
	// Determine color for each position modulo 4
	mapping := [4]byte{'?', '?', '?', '?'}
	for i := 0; i < len(s); i++ {
		if s[i] != '!' {
			mapping[i%4] = s[i]
		}
	}
	counts := [4]int{}
	for i := 0; i < len(s); i++ {
		if s[i] == '!' {
			switch mapping[i%4] {
			case 'R':
				counts[0]++
			case 'B':
				counts[1]++
			case 'Y':
				counts[2]++
			case 'G':
				counts[3]++
			}
		}
	}
	fmt.Printf("%d %d %d %d\n", counts[0], counts[1], counts[2], counts[3])
}
