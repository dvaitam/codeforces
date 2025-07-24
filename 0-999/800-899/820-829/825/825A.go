package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program decodes numbers encoded using Polycarp's binary protocol
// described in problemA.txt. Each decimal digit is represented by a sequence
// of '1' characters whose length equals the digit value, and digits are
// separated by single '0' characters. Zero digits therefore appear as two
// consecutive separators. Given the encoded string, the program prints the
// original number.
func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var s string
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}

	result := make([]byte, 0, len(s))
	count := 0
	for i := 0; i < n; i++ {
		if s[i] == '1' {
			count++
		} else { // encountered a separator
			result = append(result, byte('0'+count))
			count = 0
		}
	}
	result = append(result, byte('0'+count))

	fmt.Println(string(result))
}
