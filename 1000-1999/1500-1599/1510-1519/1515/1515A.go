package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		solve(reader, writer)
	}
}

func solve(r *bufio.Reader, w *bufio.Writer) {
	var n, x int
	fmt.Fscan(r, &n, &x)
	arr := make([]int, n)
	sum := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &arr[i])
		sum += arr[i]
	}
	if sum == x {
		fmt.Fprintln(w, "NO")
		return
	}
	sort.Ints(arr)
	prefix := 0
	for i := 0; i < n; i++ {
		prefix += arr[i]
		if prefix == x {
			// swap with the last element to avoid prefix sum == x
			arr[i], arr[n-1] = arr[n-1], arr[i]
			break
		}
	}
	fmt.Fprintln(w, "YES")
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprint(w, arr[i])
	}
	fmt.Fprintln(w)
}
