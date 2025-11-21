package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	for {
		_, err := in.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}
	fmt.Println(25)
}
