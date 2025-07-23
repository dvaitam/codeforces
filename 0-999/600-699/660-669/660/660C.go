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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &arr[i])
	}

	l, zeroes := 0, 0
	bestLen, bestL := 0, 0
	for r := 0; r < n; r++ {
		if arr[r] == 0 {
			zeroes++
		}
		for zeroes > k {
			if arr[l] == 0 {
				zeroes--
			}
			l++
		}
		if r-l+1 > bestLen {
			bestLen = r - l + 1
			bestL = l
		}
	}

	for i := bestL; i < bestL+bestLen; i++ {
		arr[i] = 1
	}

	fmt.Fprintln(writer, bestLen)
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, arr[i])
	}
}
