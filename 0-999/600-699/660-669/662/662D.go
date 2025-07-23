package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	pow10 := make([]int64, 10)
	pow10[0] = 1
	for i := 1; i < 10; i++ {
		pow10[i] = pow10[i-1] * 10
	}
	start := make([]int64, 10)
	start[1] = 1989
	for i := 2; i < 10; i++ {
		start[i] = start[i-1] + pow10[i-1]
	}

	for ; n > 0; n-- {
		var abbr string
		fmt.Fscan(reader, &abbr)
		digits := strings.TrimSpace(abbr[4:])
		k := len(digits)
		val, _ := strconv.ParseInt(digits, 10, 64)
		year := val
		mod := pow10[k]
		for year < start[k] {
			year += mod
		}
		fmt.Fprintln(writer, year)
	}
}
