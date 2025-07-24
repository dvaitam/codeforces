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

	var n int
	var w int64
	if _, err := fmt.Fscan(reader, &n, &w); err != nil {
		return
	}

	prefix := int64(0)
	minPref := int64(0)
	maxPref := int64(0)
	for i := 0; i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		prefix += x
		if prefix < minPref {
			minPref = prefix
		}
		if prefix > maxPref {
			maxPref = prefix
		}
	}

	low := int64(0)
	if -minPref > low {
		low = -minPref
	}
	high := w - maxPref
	if high > w {
		high = w
	}

	ans := int64(0)
	if high >= low {
		ans = high - low + 1
	}
	fmt.Fprintln(writer, ans)
}
