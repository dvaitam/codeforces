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
	var ansx, ansy int64
	for i := 0; i < 2*n; i++ {
		var x, y int64
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return
		}
		ansx += x
		ansy += y
	}
	fmt.Println(ansx/int64(n), ansy/int64(n))
}
