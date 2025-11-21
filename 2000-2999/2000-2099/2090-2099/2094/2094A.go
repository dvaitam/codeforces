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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s1, s2, s3 string
		fmt.Fscan(in, &s1, &s2, &s3)
		res := string([]byte{s1[0], s2[0], s3[0]})
		fmt.Fprintln(out, res)
	}
}
