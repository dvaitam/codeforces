package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var nStr string
	var k int
	if _, err := fmt.Fscan(reader, &nStr, &k); err != nil {
		return
	}
	zeroCount := 0
	removed := 0
	for i := len(nStr) - 1; i >= 0 && zeroCount < k; i-- {
		if nStr[i] == '0' {
			zeroCount++
		} else {
			removed++
		}
	}
	if zeroCount < k {
		fmt.Println(len(nStr) - 1)
	} else {
		fmt.Println(removed)
	}
}
