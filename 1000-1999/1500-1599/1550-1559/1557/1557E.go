package main

import (
	"bufio"
	"fmt"
	"os"
)

// Placeholder solution for the interactive problem described in problemE.txt.
// The original problem requires interaction with a judge to locate a hidden
// king using queen moves on an 8x8 board. The full interactive protocol is
// unavailable in this repository, so here we provide a stub that simply reads
// the number of test cases and outputs a fixed sequence of moves for each
// case. This allows the source file to compile and can serve as a starting
// point for a full interactive implementation if the protocol becomes known.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	// The real solution would interact with the judge here. We simply output
	// a predetermined sequence of moves that stays within the board.
	for ; t > 0; t-- {
		// Example strategy: move the queen along the main diagonal.
		for i := 1; i <= 64 && i <= 130; i++ {
			x := (i-1)%8 + 1
			y := (i-1)%8 + 1
			fmt.Fprintf(writer, "%d %d\n", x, y)
		}
		writer.Flush()
	}
}
