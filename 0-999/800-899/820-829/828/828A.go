package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, a, b int
	if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
		return
	}
	denied := 0
	half := 0 // two-seater tables occupied by one person
	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(in, &t)
		if t == 1 {
			if a > 0 {
				a--
			} else if b > 0 {
				b--
				half++
			} else if half > 0 {
				half--
			} else {
				denied++
			}
		} else { // t == 2
			if b > 0 {
				b--
			} else {
				denied += 2
			}
		}
	}
	fmt.Println(denied)
}
