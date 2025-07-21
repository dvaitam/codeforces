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

	var n, d int
	if _, err := fmt.Fscan(reader, &n, &d); err != nil {
		return
	}
	totalSongTime := 0
	for i := 0; i < n; i++ {
		var t int
		fmt.Fscan(reader, &t)
		totalSongTime += t
	}
	// Devu needs 10 minutes rest after each song except the last
	minRequired := totalSongTime + 10*(n-1)
	if minRequired > d {
		fmt.Fprintln(writer, -1)
		return
	}
	// During each rest period of 10 minutes, Churu can crack 2 jokes (5 min each)
	jokesDuringRest := (n - 1) * 2
	// Any remaining time can be used for extra jokes
	extraTime := d - minRequired
	extraJokes := extraTime / 5
	fmt.Fprintln(writer, jokesDuringRest+extraJokes)
}
