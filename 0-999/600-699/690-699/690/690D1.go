package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var R, C int
	if _, err := fmt.Fscan(reader, &R, &C); err != nil {
		return
	}
	cols := make([]bool, C)
	for i := 0; i < R; i++ {
		var line string
		fmt.Fscan(reader, &line)
		for j := 0; j < C && j < len(line); j++ {
			if line[j] == 'B' {
				cols[j] = true
			}
		}
	}
	seg := 0
	prev := false
	for j := 0; j < C; j++ {
		if cols[j] {
			if !prev {
				seg++
			}
			prev = true
		} else {
			prev = false
		}
	}
	fmt.Println(seg)
}
