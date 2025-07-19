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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for T > 0 {
		T--
		var n, k, r, c int
		fmt.Fscan(reader, &n, &k, &r, &c)
		// zero-index
		r--
		c--
		dif := (r%k - c%k + k) % k
		// process rows
		for i := 0; i < n; i++ {
			// build row
			row := make([]byte, n)
			im := i % k
			for j := 0; j < n; j++ {
				jm := j % k
				t := im - jm
				if t < 0 {
					t += k
				}
				if t == dif {
					row[j] = 'X'
				} else {
					row[j] = '.'
				}
			}
			writer.Write(row)
			writer.WriteByte('\n')
		}
	}
}
