package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var s string
	if _, err := fmt.Fscan(reader, &s); err != nil {
		return
	}
	pearls := 0
	links := 0
	for _, ch := range s {
		if ch == 'o' {
			pearls++
		} else if ch == '-' {
			links++
		}
	}
	if pearls <= 1 || links%pearls == 0 {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}
