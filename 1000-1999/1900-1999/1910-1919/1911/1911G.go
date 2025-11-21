package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	strs := make([][]byte, k)
	for i := 0; i < k; i++ {
		var s string
		fmt.Fscan(in, &s)
		strs[i] = []byte(s)
	}
	base := append([]byte(nil), strs[0]...)
	n := len(base)

	best := ""
	check := func(candidate []byte) {
		s := string(candidate)
		if best != "" && s >= best {
			return
		}
		for _, t := range strs {
			diff := 0
			for i := 0; i < n; i++ {
				if candidate[i] != t[i] {
					diff++
					if diff > 1 {
						break
					}
				}
			}
			if diff > 1 {
				return
			}
		}
		best = s
	}

	check(base)
	for i := 0; i < n; i++ {
		orig := base[i]
		for c := byte('a'); c <= 'z'; c++ {
			if c == orig {
				continue
			}
			base[i] = c
			check(base)
		}
		base[i] = orig
	}

	if best == "" {
		fmt.Println("-1")
	} else {
		fmt.Println(best)
	}
}
