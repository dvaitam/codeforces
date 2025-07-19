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

	var n, a, b int
	if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
		return
	}
	// Special small cases
	if n == 1 {
		fmt.Fprintln(writer, "YES")
		fmt.Fprintln(writer, "0")
		return
	}
	if (n == 2 || n == 3) && a == 1 && b == 1 {
		fmt.Fprintln(writer, "NO")
		return
	}
	// initialize matrix
	s := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			row[j] = '0'
		}
		s[i] = row
	}
	// Case: a > 1
	if a > 1 {
		if b != 1 {
			fmt.Fprintln(writer, "NO")
			return
		}
		fmt.Fprintln(writer, "YES")
		// connect first (n-a+1) nodes in a path
		val := n - a + 1
		for i := 0; i < val-1; i++ {
			s[i][i+1] = '1'
			s[i+1][i] = '1'
		}
		// output
		for i := 0; i < n; i++ {
			writer.Write(s[i])
			writer.WriteByte('\n')
		}
		return
	}
	// Case: a == 1
	// subcase: b == 1
	if b == 1 {
		fmt.Fprintln(writer, "YES")
		for i := 0; i < n-1; i++ {
			s[i][i+1] = '1'
			s[i+1][i] = '1'
		}
		for i := 0; i < n; i++ {
			writer.Write(s[i])
			writer.WriteByte('\n')
		}
		return
	}
	// subcase: b > 1
	fmt.Fprintln(writer, "YES")
	// full graph
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				s[i][j] = '1'
			}
		}
	}
	// remove edges to form b components
	val := n - b + 1
	for i := 0; i < val; i++ {
		s[0][i] = '0'
		s[i][0] = '0'
	}
	for i := 0; i < n; i++ {
		writer.Write(s[i])
		writer.WriteByte('\n')
	}
}
