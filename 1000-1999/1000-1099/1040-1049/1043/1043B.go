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
	dif := make([]int, n)
	if n > 0 {
		dif[0] = arr[0]
	}
	for i := 1; i < n; i++ {
		dif[i] = arr[i] - arr[i-1]
	}
	per := make([]bool, n+1)
	ans := n
	for i := 1; i <= n; i++ {
		per[i] = true
		for x := 0; x < n; x++ {
			if dif[x] != dif[x%i] {
				per[i] = false
				ans--
				break
			}
		}
	}
	fmt.Fprintln(writer, ans)
	for i := 1; i <= n; i++ {
		if per[i] {
			fmt.Fprintf(writer, "%d ", i)
		}
	}
	fmt.Fprintln(writer)
}
