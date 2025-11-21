package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	valid := make(map[int]struct{})
	for x := 2; ; x++ {
		val := 10
		tmp := x
		pow := 1
		for tmp > 0 {
			pow *= 10
			tmp /= 10
		}
		val = 10*pow + x
		if val > 10000 {
			break
		}
		valid[val] = struct{}{}
	}

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a int
		fmt.Fscan(in, &a)
		if _, ok := valid[a]; ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
