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

	var a, b uint64
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}

	if b-a >= 5 {
		fmt.Fprintln(writer, 0)
		return
	}

	res := uint64(1)
	for i := a + 1; i <= b; i++ {
		res = (res * (i % 10)) % 10
	}
	fmt.Fprintln(writer, res%10)
}
