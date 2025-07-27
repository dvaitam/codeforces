package main

import (
	"bufio"
	"fmt"
	"os"
)

// Since the actual problem description is missing, this program reads
// 21 lines of 21 characters (each '0' or '1') and outputs the total
// count of '1' characters.
func main() {
	in := bufio.NewReader(os.Stdin)
	count := 0
	for i := 0; i < 21; i++ {
		line, err := in.ReadString('\n')
		if err != nil && len(line) == 0 {
			return
		}
		for j := 0; j < len(line); j++ {
			if line[j] == '1' {
				count++
			}
		}
	}
	fmt.Println(count)
}
