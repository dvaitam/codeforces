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

	a := make([]int, 4)
	for i := 0; i < 4; i++ {
		fmt.Fscan(reader, &a[i])
	}
	sort.Ints(a)
	if a[0]+a[3] == a[1]+a[2] || a[0]+a[1]+a[2] == a[3] {
		fmt.Fprintln(writer, "YES")
	} else {
		fmt.Fprintln(writer, "NO")
	}
}
