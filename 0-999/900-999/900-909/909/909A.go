package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var first, last string
	if _, err := fmt.Fscan(in, &first, &last); err != nil {
		return
	}
	res := []byte{first[0]}
	for i := 1; i < len(first) && first[i] < last[0]; i++ {
		res = append(res, first[i])
	}
	res = append(res, last[0])
	fmt.Println(string(res))
}
