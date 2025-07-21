package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	// Read string (message)
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	var k int
	// Read shift key
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return
	}
	// Apply Caesar cipher
	res := make([]rune, len(s))
	for i, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			// normalize to 0-25, apply shift, wrap around
			shifted := (int(ch-'A') + k) % 26
			res[i] = rune('A' + shifted)
		} else {
			// leave other characters unchanged (though per problem all are uppercase letters)
			res[i] = ch
		}
	}
	// Output result
	fmt.Println(string(res))
}
