package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Pupil represents a pupil with name and height
type Pupil struct {
	name   string
	height int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	pupils := make([]Pupil, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &pupils[i].name, &pupils[i].height)
	}

	sort.Slice(pupils, func(i, j int) bool {
		return pupils[i].height < pupils[j].height
	})

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, pupils[i].name)
	}
}
