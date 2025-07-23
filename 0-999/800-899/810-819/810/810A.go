package main

import (
	"bufio"
	"fmt"
	"math"
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
	sum := 0
	for i := 0; i < n; i++ {
		var x int
		fmt.Fscan(reader, &x)
		sum += x
	}
	ans := 0
	for {
		avg := float64(sum) / float64(n)
		if int(math.Floor(avg+0.5)) >= k {
			break
		}
		sum += k
		n++
		ans++
	}
	fmt.Fprintln(writer, ans)
}
