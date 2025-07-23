package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int64
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	ans := int64(1)
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			ans *= i
			for n%i == 0 {
				n /= i
			}
		}
	}
	if n > 1 {
		ans *= n
	}
	fmt.Println(ans)
}
