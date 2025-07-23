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
	hasOne := false
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		if v == 1 {
			hasOne = true
		}
	}
	if hasOne {
		fmt.Println(-1)
	} else {
		fmt.Println(1)
	}
}
