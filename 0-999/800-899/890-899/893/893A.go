package main

import (
	"fmt"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	// players currently playing
	p1, p2 := 1, 2 // Alex and Bob
	spectator := 3 // Carl

	for i := 0; i < n; i++ {
		var win int
		fmt.Scan(&win)
		if win != p1 && win != p2 {
			fmt.Println("NO")
			return
		}
		// determine loser among current players
		var loser int
		if win == p1 {
			loser = p2
		} else {
			loser = p1
		}
		// next round players: winner and previous spectator
		p1, p2, spectator = win, spectator, loser
	}

	fmt.Println("YES")
}
