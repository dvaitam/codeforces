package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// discard rest of the line after n
	reader.ReadString('\n')
	line, _ := reader.ReadString('\n')
	words := strings.Fields(line)
	maxCap := 0
	for _, w := range words {
		count := 0
		for _, ch := range w {
			if ch >= 'A' && ch <= 'Z' {
				count++
			}
		}
		if count > maxCap {
			maxCap = count
		}
	}
	fmt.Println(maxCap)
}
