package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, p int
	if _, err := fmt.Fscan(reader, &n, &p); err != nil {
		return
	}
	arr := make([]int64, n)
	var total int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		total += arr[i]
	}
	var prefix int64
	var best int64
	for i := 0; i < n-1; i++ {
		prefix += arr[i]
		s1 := prefix % int64(p)
		s2 := (total - prefix) % int64(p)
		if cur := s1 + s2; cur > best {
			best = cur
		}
	}
	fmt.Fprintln(writer, best)
}
