package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	matrix := [][]int{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		row := make([]int, len(fields))
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return
			}
			row[i] = v
		}
		matrix = append(matrix, row)
	}

	n := len(matrix)
	out := bufio.NewWriter(os.Stdout)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, matrix[n-1-j][i])
		}
		if i+1 < n {
			fmt.Fprintln(out)
		}
	}
	out.Flush()
}
