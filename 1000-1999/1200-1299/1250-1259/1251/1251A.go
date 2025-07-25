package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the Codeforces problem "Broken Keyboard" (1251A).
// Each keyboard button either outputs one character (works correctly)
// or outputs two of the same character every time it's pressed
// (malfunctioning). Given the resulting string s, we need to determine
// which characters are guaranteed to correspond to working buttons.
//
// For any malfunctioning button, each time it is pressed two identical
// characters appear, so its appearances in s always occur in contiguous
// blocks of even length. Therefore, a character must be working if it
// appears in any block of odd length. We scan s, find lengths of all
// contiguous runs of the same letter and mark those with odd length.
// Finally the letters that were marked are printed in alphabetical order.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		good := make([]bool, 26)
		for i := 0; i < len(s); {
			j := i
			for j < len(s) && s[j] == s[i] {
				j++
			}
			if (j-i)%2 == 1 {
				good[s[i]-'a'] = true
			}
			i = j
		}
		for ch := 0; ch < 26; ch++ {
			if good[ch] {
				writer.WriteByte(byte('a' + ch))
			}
		}
		writer.WriteByte('\n')
	}
}
