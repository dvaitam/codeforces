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

	var layout1, layout2, s string
	if _, err := fmt.Fscan(in, &layout1, &layout2, &s); err != nil {
		return
	}

	mapping := make(map[byte]byte, 52)
	for i := 0; i < 26; i++ {
		c1 := layout1[i]
		c2 := layout2[i]
		mapping[c1] = c2
		// uppercase
		mapping[c1-'a'+'A'] = c2 - 'a' + 'A'
	}

	res := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if v, ok := mapping[ch]; ok {
			res[i] = v
		} else {
			res[i] = ch
		}
	}

	fmt.Fprintln(out, string(res))
}
