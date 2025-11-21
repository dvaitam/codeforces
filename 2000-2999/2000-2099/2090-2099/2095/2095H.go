package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var x int
	fmt.Fscan(in, &x)
	words := []string{
		"", // 1-based
		"CODEFORCES",
		"FORYOU",
		"WONT",
		"BELIEVE",
		"YOUR",
		"EYES",
		"IF",
		"YOU",
		"READ",
		"THIS",
		"LATE",
	}
	fmt.Fprintln(out, words[x])
}
