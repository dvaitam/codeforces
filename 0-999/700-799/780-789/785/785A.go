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
	total := 0
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		switch s {
		case "Tetrahedron":
			total += 4
		case "Cube":
			total += 6
		case "Octahedron":
			total += 8
		case "Dodecahedron":
			total += 12
		case "Icosahedron":
			total += 20
		}
	}
	fmt.Fprintln(writer, total)
}
