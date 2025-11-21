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
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}

	digits := make([]int, len(s))
	maxDigit := 0
	for i, ch := range s {
		digits[i] = int(ch - '0')
		if digits[i] > maxDigit {
			maxDigit = digits[i]
		}
	}

	results := make([]string, 0, maxDigit)
	for k := 0; k < maxDigit; k++ {
		buf := make([]byte, len(digits))
		for i, d := range digits {
			if d > k {
				buf[i] = '1'
			} else {
				buf[i] = '0'
			}
		}
		// trim leading zeros
		idx := 0
		for idx < len(buf) && buf[idx] == '0' {
			idx++
		}
		if idx == len(buf) {
			continue
		}
		results = append(results, string(buf[idx:]))
	}

	if len(results) == 0 {
		results = append(results, "0")
	}

	fmt.Fprintln(out, len(results))
	for i, val := range results {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, val)
	}
	fmt.Fprintln(out)
}
