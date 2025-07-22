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

	var v1, v2, v3, vm int
	if _, err := fmt.Fscan(reader, &v1, &v2, &v3, &vm); err != nil {
		return
	}

	for a := 2 * v1; a >= v1; a-- {
		if a <= 2*vm { // Masha must dislike this car, so a must be > 2*vm
			continue
		}
		if vm > a { // Masha must be able to climb into it
			continue
		}
		for b := min(2*v2, a-1); b >= v2; b-- {
			if b <= 2*vm {
				continue
			}
			if vm > b {
				continue
			}
			for c := min(2*v3, b-1); c >= v3; c-- {
				if vm > c {
					continue
				}
				if 2*vm < c { // Masha must like only smallest car -> likes if 2*vm >= c
					continue
				}
				fmt.Fprintf(writer, "%d %d %d\n", a, b, c)
				return
			}
		}
	}

	fmt.Fprintln(writer, "-1")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
