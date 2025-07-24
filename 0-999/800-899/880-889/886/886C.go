package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	t := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &t[i])
	}

	rooms := map[int]struct{}{0: {}}

	for i := 1; i <= n; i++ {
		prev := t[i]
		if _, ok := rooms[prev]; ok {
			delete(rooms, prev)
		}
		rooms[i] = struct{}{}
	}

	fmt.Println(len(rooms))
}
