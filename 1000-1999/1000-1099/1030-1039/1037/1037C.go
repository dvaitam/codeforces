package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	var aStr, bStr string
	fmt.Fscan(in, &aStr, &bStr)
	a := []byte(aStr)
	b := []byte(bStr)
	ans := 0
	for i := 1; i < n; i++ {
		if a[i-1] == b[i] && a[i] == b[i-1] && a[i-1] != a[i] {
			ans++
			a[i-1] = b[i-1]
			a[i] = b[i]
		}
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			ans++
		}
	}
	fmt.Println(ans)
}
