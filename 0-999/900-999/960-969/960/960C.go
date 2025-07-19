package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var x, d int64
	if _, err := fmt.Fscan(reader, &x, &d); err != nil {
		return
	}
	ans := make([]int64, 0)
	last := int64(1)
	// For small x, output directly
	if x <= 10000 {
		fmt.Println(x)
		for i := int64(0); i < x; i++ {
			fmt.Printf("%d ", last)
			last += d + 2
		}
		fmt.Println()
		return
	}
	// Build answer using binary decomposition
	for x > 0 {
		// find largest len such that (1<<len)-1 <= x
		l := 1
		for {
			v := (1 << l) - 1
			if int64(v) > x {
				break
			}
			l++
		}
		l--
		// append l copies of last
		for i := 0; i < l; i++ {
			ans = append(ans, last)
		}
		last += d + 2
		x -= int64((1 << l) - 1)
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, len(ans))
	for i, v := range ans {
		if i > 0 {
			writer.WriteByte(' ')
		}
		writer.WriteString(fmt.Sprintf("%d", v))
	}
	writer.WriteByte('\n')
}
