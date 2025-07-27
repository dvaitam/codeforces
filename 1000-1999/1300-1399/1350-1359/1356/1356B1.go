package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	bits := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &bits[i])
	}
	// increment with carry
	carry := 1
	for i := 0; i < n; i++ {
		if carry == 0 {
			break
		}
		if bits[i] == 0 {
			bits[i] = 1
			carry = 0
		} else {
			bits[i] = 0
		}
	}
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < n; i++ {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, bits[i])
	}
	writer.WriteByte('\n')
}
