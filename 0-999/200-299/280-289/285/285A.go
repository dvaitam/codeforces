package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var N, K int
	if _, err := fmt.Fscan(os.Stdin, &N, &K); err != nil {
		return
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	// print from N down to N-K+1
	for i := N; i > N-K; i-- {
		fmt.Fprintf(w, "%d ", i)
	}
	// print from 1 to N-K
	for j := 1; j <= N-K; j++ {
		fmt.Fprintf(w, "%d ", j)
	}
	fmt.Fprintln(w)
}
