package main

import (
	"bufio"
	"fmt"
	"os"
)

func build(n int) []int {
	if n == 2 || n == 3 {
		return nil
	}
	if n == 4 {
		return []int{2, 4, 1, 3}
	}
	res := make([]int, 0, n)
	// append odd numbers ascending
	for i := 1; i <= n; i += 2 {
		res = append(res, i)
	}
	if n%2 == 0 {
		evenSeq := []int{n - 4, n, n - 2}
		for x := n - 6; x >= 2; x -= 2 {
			evenSeq = append(evenSeq, x)
		}
		res = append(res, evenSeq...)
	} else {
		evenSeq := []int{n - 3, n - 1}
		for x := n - 5; x >= 2; x -= 2 {
			evenSeq = append(evenSeq, x)
		}
		res = append(res, evenSeq...)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		p := build(n)
		if p == nil {
			fmt.Fprintln(writer, -1)
			continue
		}
		for i, v := range p {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, v)
		}
		fmt.Fprintln(writer)
	}
}
