package main

import (
	"fmt"
	"os"
)

func main() {
	var s string
	if err := fmt.Fscan(os.Stdin, &s); err != nil {
		return
	}
	res, num, op := 0, 0, 1
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
		} else {
			res += num * op
			num = 0
			if ch == '+' {
				op = 1
			} else {
				op = -1
			}
		}
	}
	res += num * op

	// hundreds digit
	if res >= 100 {
		for i := 0; i < res/100; i++ {
			fmt.Fprint(os.Stdout, "+")
		}
		fmt.Fprint(os.Stdout, "\n++++++++++++++++++++++++++++++++++++++++++++++++.")
		fmt.Fprint(os.Stdout, ">")
	}
	// tens digit
	if res >= 10 {
		for i := 0; i < (res%100)/10; i++ {
			fmt.Fprint(os.Stdout, "+")
		}
		fmt.Fprint(os.Stdout, "\n++++++++++++++++++++++++++++++++++++++++++++++++.")
		fmt.Fprint(os.Stdout, ">")
	}
	// ones digit
	for i := 0; i < res%10; i++ {
		fmt.Fprint(os.Stdout, "+")
	}
	fmt.Fprint(os.Stdout, "\n++++++++++++++++++++++++++++++++++++++++++++++++.")
}
