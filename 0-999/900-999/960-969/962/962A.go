package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]int, n)
	total := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		total += arr[i]
	}
	need := (total + 1) / 2
	sum := 0
	for i, v := range arr {
		sum += v
		if sum >= need {
			fmt.Println(i + 1)
			return
		}
	}
}
