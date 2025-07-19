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
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)
	fmt.Fprintf(writer, "1 %d\n", a[0])
	if a[n-1] > 0 {
		fmt.Fprintf(writer, "1 %d\n", a[n-1])
		fmt.Fprintf(writer, "%d", n-2)
		for i := 1; i < n-1; i++ {
			fmt.Fprintf(writer, " %d", a[i])
		}
		fmt.Fprint(writer, "\n")
	} else {
		fmt.Fprintf(writer, "2 %d %d\n", a[1], a[2])
		fmt.Fprintf(writer, "%d", n-3)
		for i := 3; i < n; i++ {
			fmt.Fprintf(writer, " %d", a[i])
		}
		fmt.Fprint(writer, "\n")
	}
}
