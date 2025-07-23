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

	var aStr, bStr string
	if _, err := fmt.Fscan(reader, &aStr); err != nil {
		return
	}
	if _, err := fmt.Fscan(reader, &bStr); err != nil {
		return
	}

	if len(aStr) < len(bStr) {
		digits := []byte(aStr)
		sort.Slice(digits, func(i, j int) bool { return digits[i] > digits[j] })
		fmt.Fprintln(writer, string(digits))
		return
	}

	counts := make([]int, 10)
	for _, ch := range aStr {
		counts[ch-'0']++
	}
	res := make([]byte, len(bStr))
	if dfs(0, true, counts, []byte(bStr), res) {
		fmt.Fprintln(writer, string(res))
	}
}

func dfs(pos int, limit bool, counts []int, bDigits []byte, res []byte) bool {
	if pos == len(bDigits) {
		return true
	}
	maxDigit := 9
	if limit {
		maxDigit = int(bDigits[pos] - '0')
	}
	for d := maxDigit; d >= 0; d-- {
		if counts[d] == 0 {
			continue
		}
		counts[d]--
		res[pos] = byte('0' + d)
		if dfs(pos+1, limit && d == maxDigit, counts, bDigits, res) {
			return true
		}
		counts[d]++
	}
	return false
}
