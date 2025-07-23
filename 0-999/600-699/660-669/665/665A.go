package main

import (
	"bufio"
	"fmt"
	"os"
)

func parseTime(s string) int {
	var h, m int
	fmt.Sscanf(s, "%d:%d", &h, &m)
	return h*60 + m
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var a, ta int
	if _, err := fmt.Fscan(reader, &a, &ta); err != nil {
		return
	}
	_ = a
	var b, tb int
	if _, err := fmt.Fscan(reader, &b, &tb); err != nil {
		return
	}
	var dep string
	if _, err := fmt.Fscan(reader, &dep); err != nil {
		return
	}

	start := parseTime(dep)
	end := start + ta

	const dayStart = 5 * 60
	const dayEnd = 23*60 + 59

	count := 0
	for t := dayStart; t <= dayEnd; t += b {
		// interval of a bus from B to A
		tEnd := t + tb
		if t < end && tEnd > start && t != end && tEnd != start {
			count++
		}
	}

	fmt.Fprintln(writer, count)
}
