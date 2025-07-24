package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	pattern := "Yes"
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		ok := false
		for start := 0; start < 3; start++ {
			good := true
			for i := 0; i < len(s); i++ {
				if s[i] != pattern[(start+i)%3] {
					good = false
					break
				}
			}
			if good {
				ok = true
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
