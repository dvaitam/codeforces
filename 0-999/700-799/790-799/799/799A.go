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

	var n, t, k, d int
	if _, err := fmt.Fscan(reader, &n, &t, &k, &d); err != nil {
		return
	}

	// Time needed using only one oven
	batches := (n + k - 1) / k
	oneTime := batches * t

	for time := 1; time < oneTime; time++ {
		cakes := (time / t) * k
		if time > d {
			cakes += ((time - d) / t) * k
		}
		if cakes >= n {
			fmt.Fprintln(writer, "YES")
			return
		}
	}
	fmt.Fprintln(writer, "NO")
}
