package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var s string
	fmt.Fscan(reader, &s)

	n := new(big.Int)
	n.SetString(s, 2)

	pow := big.NewInt(1)
	count := 0
	for pow.Cmp(n) < 0 {
		count++
		pow.Lsh(pow, 2)
	}

	fmt.Fprintln(writer, count)
}
