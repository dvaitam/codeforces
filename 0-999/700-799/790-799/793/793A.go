package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	var k int64
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int64, n)
	var minA int64 = 1<<63 - 1
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
		if arr[i] < minA {
			minA = arr[i]
		}
	}
	remainder := arr[0] % k
	for i := 1; i < n; i++ {
		if arr[i]%k != remainder {
			fmt.Println(-1)
			return
		}
	}
	var ans int64
	for i := 0; i < n; i++ {
		ans += (arr[i] - minA) / k
	}
	fmt.Println(ans)
}
