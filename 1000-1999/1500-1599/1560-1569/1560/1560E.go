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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var t string
		fmt.Fscan(in, &t)

		// Determine removal order
		seen := make(map[byte]bool)
		orderRev := make([]byte, 0)
		for i := len(t) - 1; i >= 0; i-- {
			ch := t[i]
			if !seen[ch] {
				seen[ch] = true
				orderRev = append(orderRev, ch)
			}
		}
		order := make([]byte, len(orderRev))
		for i := range orderRev {
			order[len(orderRev)-1-i] = orderRev[i]
		}

		// frequency of all characters in t
		freq := make(map[byte]int)
		for i := 0; i < len(t); i++ {
			freq[t[i]]++
		}

		origCount := make(map[byte]int)
		prefLen := 0
		valid := true
		for i, ch := range order {
			c := freq[ch]
			step := i + 1
			if c%step != 0 {
				valid = false
				break
			}
			origCount[ch] = c / step
			prefLen += origCount[ch]
		}
		if !valid {
			fmt.Fprintln(out, -1)
			continue
		}
		if prefLen > len(t) {
			fmt.Fprintln(out, -1)
			continue
		}

		s := t[:prefLen]

		// simulate process
		cur := []byte(s)
		built := make([]byte, 0, len(t))
		for _, ch := range order {
			built = append(built, cur...)
			filtered := make([]byte, 0, len(cur))
			for _, c := range cur {
				if c != ch {
					filtered = append(filtered, c)
				}
			}
			cur = filtered
		}
		if string(built) != t {
			fmt.Fprintln(out, -1)
			continue
		}
		fmt.Fprintf(out, "%s %s\n", s, string(order))
	}
}
