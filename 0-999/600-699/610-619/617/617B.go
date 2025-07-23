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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	positions := make([]int, 0)
	for i, v := range arr {
		if v == 1 {
			positions = append(positions, i)
		}
	}

	if len(positions) == 0 {
		fmt.Fprintln(writer, 0)
		return
	}

	ans := int64(1)
	for i := 1; i < len(positions); i++ {
		diff := positions[i] - positions[i-1]
		ans *= int64(diff)
	}

	fmt.Fprintln(writer, ans)
}
