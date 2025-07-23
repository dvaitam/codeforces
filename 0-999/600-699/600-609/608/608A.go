package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, s int
	if _, err := fmt.Fscan(reader, &n, &s); err != nil {
		return
	}
	arrival := make([]int, s+1)
	for i := 0; i < n; i++ {
		var f, t int
		fmt.Fscan(reader, &f, &t)
		if arrival[f] < t {
			arrival[f] = t
		}
	}
	time := 0
	for floor := s; floor > 0; floor-- {
		if arrival[floor] > time {
			time = arrival[floor]
		}
		time++
	}
	fmt.Println(time)
}
