package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s string
	fmt.Fscan(in, &s)

	var ans int64
	for i, ch := range s {
		d := int(ch - '0')
		if d%4 == 0 {
			ans++
		}
		if i > 0 {
			prev := int(s[i-1]-'0')*10 + d
			if prev%4 == 0 {
				ans += int64(i)
			}
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprint(out, ans)
	out.Flush()
}
