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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	horiz := make([]bool, n+1)
	vert := make([]bool, n+1)

	res := make([]int, 0, n*n)
	for day := 1; day <= n*n; day++ {
		var h, v int
		fmt.Fscan(reader, &h, &v)
		if !horiz[h] && !vert[v] {
			horiz[h] = true
			vert[v] = true
			res = append(res, day)
		}
	}

	for i, d := range res {
		if i > 0 {
			writer.WriteByte(' ')
		}
		fmt.Fprint(writer, d)
	}
	if len(res) > 0 {
		writer.WriteByte('\n')
	}
}
