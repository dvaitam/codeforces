package main

import (
	"bufio"
	"fmt"
	"os"
)

func query(nums []int, in *bufio.Reader, out *bufio.Writer) int {
	fmt.Fprint(out, "?")
	for _, v := range nums {
		fmt.Fprint(out, " ", v)
	}
	fmt.Fprintln(out)
	out.Flush()
	var resp int
	fmt.Fscan(in, &resp)
	return resp
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q1 []int
	for i := 1; i <= 100; i++ {
		q1 = append(q1, i)
	}
	r1 := query(q1, in, out)

	var q2 []int
	for i := 1; i <= 100; i++ {
		q2 = append(q2, i<<7)
	}
	r2 := query(q2, in, out)

	high := r1 >> 7
	low := r2 & 127
	ans := (high << 7) | low
	fmt.Fprintln(out, "!", ans)
	out.Flush()
}
