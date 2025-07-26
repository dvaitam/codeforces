package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		i := int(s[0] - 'a' + 1)
		j := int(s[1] - 'a' + 1)
		ans := (i-1)*25 + j
		if j > i {
			ans--
		}
		fmt.Println(ans)
	}
}
