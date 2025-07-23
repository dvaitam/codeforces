package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var s, x int64
	if _, err := fmt.Fscan(in, &s, &x); err != nil {
		return
	}
	if s < x || (s-x)%2 != 0 {
		fmt.Println(0)
		return
	}
	t := (s - x) / 2
	if t&x != 0 {
		fmt.Println(0)
		return
	}
	count := 0
	tmp := x
	for tmp > 0 {
		if tmp&1 == 1 {
			count++
		}
		tmp >>= 1
	}
	result := int64(1) << count
	if t == 0 {
		result -= 2
	}
	fmt.Println(result)
}
