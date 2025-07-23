package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for t := 0; t < n; t++ {
		ok := true
		for i := 0; i < n; i++ {
			val := a[i]
			if i%2 == 0 {
				val = (val + t) % n
			} else {
				val = ((val-t)%n + n) % n
			}
			if val != i {
				ok = false
				break
			}
		}
		if ok {
			fmt.Println("Yes")
			return
		}
	}
	fmt.Println("No")
}
