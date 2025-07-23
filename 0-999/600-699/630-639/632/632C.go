package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	arr := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i]+arr[j] < arr[j]+arr[i]
	})
	// Concatenate result
	var b strings.Builder
	for _, s := range arr {
		b.WriteString(s)
	}
	fmt.Println(b.String())
}
