package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// This program follows the interaction protocol described in problemA2.txt.
// It makes up to 100 attempts, printing randomly generated bit strings of
// length 5000 and reading the result after each attempt.
func main() {
	const n = 5000
	const attempts = 100

	rand.Seed(time.Now().UnixNano())

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	best := 0
	for i := 0; i < attempts; i++ {
		// Generate random answers for this attempt.
		buf := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				buf[j] = '0'
			} else {
				buf[j] = '1'
			}
		}
		// Print the answers and flush to ensure they are sent immediately.
		fmt.Fprintln(out, string(buf))
		out.Flush()

		var result int
		if _, err := fmt.Fscan(in, &result); err != nil {
			return
		}
		if result > best {
			best = result
		}
		// If the exam was passed (result == n+1), no more attempts are needed.
		if result == n+1 {
			break
		}
	}
	_ = best // best can be used for debugging or scoring if needed
}
