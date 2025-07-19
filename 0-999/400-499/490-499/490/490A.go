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

	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return
	}
	team1 := make([]int, 0, m)
	team2 := make([]int, 0, m)
	team3 := make([]int, 0, m)
	for i := 1; i <= m; i++ {
		var t int
		fmt.Fscan(reader, &t)
		switch t {
		case 1:
			team1 = append(team1, i)
		case 2:
			team2 = append(team2, i)
		case 3:
			team3 = append(team3, i)
		}
	}
	minx := len(team1)
	if len(team2) < minx {
		minx = len(team2)
	}
	if len(team3) < minx {
		minx = len(team3)
	}
	fmt.Fprintln(writer, minx)
	for i := 0; i < minx; i++ {
		fmt.Fprintf(writer, "%d %d %d\n", team1[i], team2[i], team3[i])
	}
}
