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

	var A, B int64
	if _, err := fmt.Fscan(in, &A, &B); err != nil {
		return
	}
	var x, y, z int64
	fmt.Fscan(in, &x, &y, &z)

	needYellow := 2*x + y
	needBlue := y + 3*z

	add := int64(0)
	if needYellow > A {
		add += needYellow - A
	}
	if needBlue > B {
		add += needBlue - B
	}

	fmt.Fprintln(out, add)
}
