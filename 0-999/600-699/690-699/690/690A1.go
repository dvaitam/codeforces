package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	ans := (n + 1) / 2
	fmt.Println(ans)
}
