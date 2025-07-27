package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func query(writer *bufio.Writer, reader *bufio.Reader, op string, i, j int) int {
	fmt.Fprintf(writer, "%s %d %d\n", op, i, j)
	writer.Flush()
	var x int
	fmt.Fscan(reader, &x)
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	and12 := query(writer, reader, "and", 1, 2)
	and13 := query(writer, reader, "and", 1, 3)
	and23 := query(writer, reader, "and", 2, 3)

	or12 := query(writer, reader, "or", 1, 2)
	or13 := query(writer, reader, "or", 1, 3)
	or23 := query(writer, reader, "or", 2, 3)

	sum12 := and12 + or12
	sum13 := and13 + or13
	sum23 := and23 + or23

	a := make([]int, n+1)
	a[1] = (sum12 + sum13 - sum23) / 2
	a[2] = sum12 - a[1]
	a[3] = sum13 - a[1]

	for i := 4; i <= n; i++ {
		and1i := query(writer, reader, "and", 1, i)
		or1i := query(writer, reader, "or", 1, i)
		sum1i := and1i + or1i
		a[i] = sum1i - a[1]
	}

	values := append([]int(nil), a[1:]...)
	sort.Ints(values)
	kth := values[k-1]
	fmt.Fprintf(writer, "finish %d\n", kth)
}
