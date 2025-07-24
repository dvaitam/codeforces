package main

import (
	"bufio"
	"fmt"
	"os"
)

const p = 64
const m = 512 / p

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	names := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			names = append(names, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return
	}
	permutation := make([]int, m*m)
	for i := 0; i < len(permutation); i++ {
		permutation[i] = i
	}
	for _, name := range names {
		fmt.Print(name)
		for _, v := range permutation {
			fmt.Printf(" %d", v)
		}
		fmt.Println()
	}
}
