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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		seqs := make([][]int, n)
		freq := make(map[int]int)
		for i := 0; i < n; i++ {
			seqs[i] = make([]int, n-1)
			for j := 0; j < n-1; j++ {
				fmt.Fscan(reader, &seqs[i][j])
			}
			freq[seqs[i][0]]++
		}
		// find element appearing n-1 times as the first element
		var firstVal int
		for val, cnt := range freq {
			if cnt > freq[firstVal] {
				firstVal = val
			}
		}
		// find the sequence that does not start with firstVal
		var rest []int
		for i := 0; i < n; i++ {
			if seqs[i][0] != firstVal {
				rest = seqs[i]
				break
			}
		}
		fmt.Fprint(writer, firstVal)
		for _, v := range rest {
			fmt.Fprint(writer, " ", v)
		}
		fmt.Fprint(writer, "\n")
	}
}
