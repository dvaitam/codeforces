package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program provides a minimal stub for the interactive problem described in
// problemI.txt. The real interactive protocol is not implemented here. Instead
// it simply prints a placeholder circle description and exits so the file
// compiles and can be used as a template for further development.

func main() {
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, "0 0 1")
	out.Flush()
}
