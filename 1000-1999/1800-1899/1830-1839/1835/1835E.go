package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	_ = bufio.NewReader(os.Stdin) // placeholder to match style
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	// TODO: implement solution for Codeforces problem described in problemE.txt
	fmt.Fprintln(writer, "TODO: solution not yet implemented")
}
