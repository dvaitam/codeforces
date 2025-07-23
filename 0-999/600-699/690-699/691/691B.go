package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSpace(s)

	mirror := map[rune]rune{
		'A': 'A', 'H': 'H', 'I': 'I', 'M': 'M',
		'O': 'O', 'T': 'T', 'U': 'U', 'V': 'V',
		'W': 'W', 'X': 'X', 'Y': 'Y',
		'b': 'd', 'd': 'b', 'p': 'q', 'q': 'p',
		'o': 'o', 'v': 'v', 'w': 'w', 'x': 'x',
	}

	n := len(s)
	for i := 0; i < n; i++ {
		c := rune(s[i])
		m, ok := mirror[c]
		if !ok || m != rune(s[n-1-i]) {
			fmt.Println("NIE")
			return
		}
	}
	fmt.Println("TAK")
}
