package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		solveC(in, out)
	}
}

func solveC(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Ints(a)
	if n == 1 {
		fmt.Fprintln(out, a[0])
		return
	}
	if a[n-1] == a[n-2] {
		fmt.Fprintln(out, a[n-1])
	} else {
		fmt.Fprintln(out, a[n-2])
	}
}
