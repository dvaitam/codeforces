package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var hex string
	fmt.Fscan(reader, &hex)
	var value uint64
	fmt.Sscanf(hex, "%x", &value)
	fmt.Println(value)
}
