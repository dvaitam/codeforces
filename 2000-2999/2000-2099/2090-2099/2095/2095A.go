package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// April Fools: the answer is already in the picture.
	in := bufio.NewReader(os.Stdin)
	// consume all input (if any) to follow usual pattern, though it is unused
	for {
		_, err := in.ReadByte()
		if err != nil {
			break
		}
	}
	fmt.Println("printf(\"puzzling\");")
}
