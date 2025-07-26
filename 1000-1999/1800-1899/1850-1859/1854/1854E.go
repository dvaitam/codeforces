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

	var m int64
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}

	cnt := make([]int64, 60)
	cnt[0] = 1
	ans := make([]int, 0)
	for i, r, c := 1, int64(0), 0; i <= 60 && c <= 60; {
		if r < (int64(1)<<uint(i)) && m >= cnt[60-i] {
			m -= cnt[60-i]
			for j := 59; j >= i; j-- {
				cnt[j] += cnt[j-i]
			}
			r++
			c++
			ans = append(ans, i)
		} else {
			r = 0
			i++
		}
	}

	fmt.Fprintln(writer, len(ans))
	for idx, v := range ans {
		if idx > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, v)
	}
}
