package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var a int
	if _, err := fmt.Fscan(in, &a); err != nil {
		return
	}
	res := ((a>>4)&1)<<0 |
		((a>>1)&1)<<1 |
		((a>>3)&1)<<2 |
		((a>>2)&1)<<3 |
		((a>>0)&1)<<4 |
		((a>>5)&1)<<5
	fmt.Println(res)
}
