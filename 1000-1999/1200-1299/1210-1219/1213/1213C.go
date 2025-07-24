package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n, m int64
		fmt.Fscan(in, &n, &m)
		cnt := n / m
		last := m % 10
		var pref [10]int64
		sumCycle := int64(0)
		for i := 1; i <= 10; i++ {
			digit := (last * int64(i)) % 10
			sumCycle += digit
			pref[i-1] = sumCycle
		}
		res := (cnt / 10) * sumCycle
		if r := cnt % 10; r > 0 {
			res += pref[r-1]
		}
		fmt.Fprintln(out, res)
	}
}
