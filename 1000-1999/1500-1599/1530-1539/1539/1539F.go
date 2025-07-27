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

	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	leftLess := make([]int, n+2)
	rightLess := make([]int, n+2)
	leftGreater := make([]int, n+2)
	rightGreater := make([]int, n+2)

	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			leftLess[i] = stack[len(stack)-1]
		} else {
			leftLess[i] = 0
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			rightLess[i] = stack[len(stack)-1]
		} else {
			rightLess[i] = n + 1
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			leftGreater[i] = stack[len(stack)-1]
		} else {
			leftGreater[i] = 0
		}
		stack = append(stack, i)
	}

	stack = stack[:0]
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] <= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 {
			rightGreater[i] = stack[len(stack)-1]
		} else {
			rightGreater[i] = n + 1
		}
		stack = append(stack, i)
	}

	res := make([]int, n+1)
	for i := 1; i <= n; i++ {
		lenMin := rightLess[i] - leftLess[i] - 1
		distMin := lenMin / 2
		lenMax := rightGreater[i] - leftGreater[i] - 1
		distMax := (lenMax+1)/2 - 1
		if distMin > distMax {
			res[i] = distMin
		} else {
			res[i] = distMax
		}
	}

	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, res[i])
	}
	fmt.Fprintln(out)
}
