package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var w, m int64
	if _, err := fmt.Fscan(in, &w, &m); err != nil {
		return
	}
	for m > 0 {
		r := m % w
		if r == 0 || r == 1 {
			m /= w
		} else if r == w-1 {
			m = m/w + 1
		} else {
			fmt.Println("NO")
			return
		}
	}
	fmt.Println("YES")
}
