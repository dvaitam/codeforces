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
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(reader, &n, &k)
		if k == 1 {
			fmt.Fprintln(writer, n)
			continue
		}
		ans := int64(0)
		for n > 0 {
			ans += n % k
			n /= k
		}
		fmt.Fprintln(writer, ans)
	}
}
