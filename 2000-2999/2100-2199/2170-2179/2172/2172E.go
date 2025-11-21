package main

import (
	"bufio"
	"fmt"
	"os"
)

var factorial = []int{1, 1, 2, 6, 24}

func kthPermutation(length, idx int) []byte {
	available := make([]byte, length)
	for i := 0; i < length; i++ {
		available[i] = byte('1' + i)
	}

	idx--
	result := make([]byte, length)
	for pos := 0; pos < length; pos++ {
		blockSize := factorial[length-1-pos]
		pick := idx / blockSize
		idx %= blockSize
		result[pos] = available[pick]
		available = append(available[:pick], available[pick+1:]...)
	}
	return result
}

func countMatches(a, b []byte) (int, int) {
	same := 0
	for i := range a {
		if a[i] == b[i] {
			same++
		}
	}
	different := len(a) - same
	return same, different
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, j, k int
		fmt.Fscan(in, &n, &j, &k)
		length := len(fmt.Sprint(n))
		permJ := kthPermutation(length, j)
		permK := kthPermutation(length, k)
		a, b := countMatches(permJ, permK)
		fmt.Fprintf(out, "%dA%dB\n", a, b)
	}
}
