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

	var n, l, k int
	if _, err := fmt.Fscan(reader, &n, &l, &k); err != nil {
		return
	}
	var letters string
	fmt.Fscan(reader, &letters)

	arr := []rune(letters)
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

	words := make([][]rune, n)
	pos := 0
	L, R := 0, k-1
	for i := 0; i < l; i++ {
		for j := L; j <= R; j++ {
			words[j] = append(words[j], arr[pos])
			pos++
		}
		for L < R && words[L][i] != words[R][i] {
			L++
		}
	}

	for i := 0; i < n; i++ {
		for len(words[i]) < l {
			words[i] = append(words[i], arr[pos])
			pos++
		}
	}

	res := make([]string, n)
	for i := range words {
		res[i] = string(words[i])
	}
	sort.Strings(res)
	for _, w := range res {
		fmt.Fprintln(writer, w)
	}
}
