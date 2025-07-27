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

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		rows := make([]int, n)
		cols := make([]int, m)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var x int
				fmt.Fscan(reader, &x)
				if x == 1 {
					rows[i] = 1
					cols[j] = 1
				}
			}
		}
		emptyRows := 0
		for i := 0; i < n; i++ {
			if rows[i] == 0 {
				emptyRows++
			}
		}
		emptyCols := 0
		for j := 0; j < m; j++ {
			if cols[j] == 0 {
				emptyCols++
			}
		}
		moves := emptyRows
		if emptyCols < moves {
			moves = emptyCols
		}
		if moves%2 == 1 {
			fmt.Fprintln(writer, "Ashish")
		} else {
			fmt.Fprintln(writer, "Vivek")
		}
	}
}
