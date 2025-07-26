package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var k int
	fmt.Fscan(reader, &k)

	var s, t string
	fmt.Fscan(reader, &s)
	fmt.Fscan(reader, &t)

	// TODO: implement full solution for k up to 12
	fmt.Fprintln(writer, 0)
}
