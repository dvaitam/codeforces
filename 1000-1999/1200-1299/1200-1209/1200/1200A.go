package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	rooms := make([]byte, 10)
	for i := 0; i < 10; i++ {
		rooms[i] = '0'
	}
	for _, ch := range s {
		switch ch {
		case 'L':
			for i := 0; i < 10; i++ {
				if rooms[i] == '0' {
					rooms[i] = '1'
					break
				}
			}
		case 'R':
			for i := 9; i >= 0; i-- {
				if rooms[i] == '0' {
					rooms[i] = '1'
					break
				}
			}
		default:
			rooms[ch-'0'] = '0'
		}
	}
	out.Write(rooms)
}
