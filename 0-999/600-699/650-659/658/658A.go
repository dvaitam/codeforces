package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program solves the problem described in problemA.txt (Bear and Reverse Radewoosh).
// It computes the scores for Limak and Radewoosh given the time decay and determines the winner.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, c int
	if _, err := fmt.Fscan(reader, &n, &c); err != nil {
		return
	}

	p := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	t := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &t[i])
	}

	time := 0
	scoreL := 0
	for i := 0; i < n; i++ {
		time += t[i]
		s := p[i] - c*time
		if s < 0 {
			s = 0
		}
		scoreL += s
	}

	time = 0
	scoreR := 0
	for i := n - 1; i >= 0; i-- {
		time += t[i]
		s := p[i] - c*time
		if s < 0 {
			s = 0
		}
		scoreR += s
	}

	if scoreL > scoreR {
		fmt.Fprintln(writer, "Limak")
	} else if scoreL < scoreR {
		fmt.Fprintln(writer, "Radewoosh")
	} else {
		fmt.Fprintln(writer, "Tie")
	}
}
