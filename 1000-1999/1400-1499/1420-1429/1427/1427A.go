package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	fmt.Fscan(reader, &q)
	for i := 0; i < q; i++ {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n)
		s := 0
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &a[j])
			s += a[j]
		}
		if s == 0 {
			writer.WriteString("NO\n")
			continue
		}
		writer.WriteString("YES\n")
		if s > 0 {
			for j := 0; j < n; j++ {
				if a[j] > 0 {
					writer.WriteString(strconv.Itoa(a[j]))
					writer.WriteByte(' ')
				}
			}
			for j := 0; j < n; j++ {
				if a[j] <= 0 {
					writer.WriteString(strconv.Itoa(a[j]))
					writer.WriteByte(' ')
				}
			}
			writer.WriteByte('\n')
		} else {
			for j := 0; j < n; j++ {
				if a[j] < 0 {
					writer.WriteString(strconv.Itoa(a[j]))
					writer.WriteByte(' ')
				}
			}
			for j := 0; j < n; j++ {
				if a[j] >= 0 {
					writer.WriteString(strconv.Itoa(a[j]))
					writer.WriteByte(' ')
				}
			}
			writer.WriteByte('\n')
		}
	}
}
