package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type player struct {
	val int
	id  int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	players := make([]player, n)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		players[i] = player{val: v, id: i + 1}
	}

	sort.Slice(players, func(i, j int) bool {
		return players[i].val > players[j].val
	})

	team1 := make([]int, 0, (n+1)/2)
	team2 := make([]int, 0, n/2)
	for i, p := range players {
		if i%2 == 0 {
			team1 = append(team1, p.id)
		} else {
			team2 = append(team2, p.id)
		}
	}

	fmt.Fprintln(writer, len(team1))
	for _, id := range team1 {
		fmt.Fprint(writer, id, " ")
	}
	fmt.Fprintln(writer)

	fmt.Fprintln(writer, len(team2))
	for _, id := range team2 {
		fmt.Fprint(writer, id, " ")
	}
	fmt.Fprintln(writer)
}
