package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	sum := 0
	for {
		var x int
		_, err := fmt.Fscan(reader, &x)
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		sum += x
	}
	fmt.Println(sum)
}
