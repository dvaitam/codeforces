package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	var t int
	fmt.Sscanf(in.Text(), "%d", &t)
	for i := 0; i < t; i++ {
		in.Scan()
		s := strings.TrimSpace(in.Text())
		r := []rune(s)
		ok := true
		for l, rp := 0, len(r)-1; l < rp; l, rp = l+1, rp-1 {
			if r[l] != r[rp] {
				ok = false
				break
			}
		}
		if ok {
			fmt.Println("Yes")
		} else {
			fmt.Println("No")
		}
	}
}
