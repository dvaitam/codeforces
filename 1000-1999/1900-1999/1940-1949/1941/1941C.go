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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		ans := 0
		for i := 0; i < n; {
			if i+4 < n && s[i:i+5] == "mapie" {
				ans++
				i += 5
			} else if i+2 < n && (s[i:i+3] == "pie" || s[i:i+3] == "map") {
				ans++
				i += 3
			} else {
				i++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
