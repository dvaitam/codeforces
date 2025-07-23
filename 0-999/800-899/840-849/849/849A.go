package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	if n%2 == 1 && arr[0]%2 == 1 && arr[n-1]%2 == 1 {
		fmt.Println("Yes")
	} else {
		fmt.Println("No")
	}
}
