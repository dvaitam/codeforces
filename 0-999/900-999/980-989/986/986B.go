package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &p[i])
	}
	visited := make([]bool, n+1)
	cycles := 0
	for i := 1; i <= n; i++ {
		if !visited[i] {
			cycles++
			for j := i; !visited[j]; j = p[j] {
				visited[j] = true
			}
		}
	}
	parity := (n - cycles) % 2
	petrParity := n % 2
	if parity == petrParity {
		fmt.Println("Petr")
	} else {
		fmt.Println("Um_nik")
	}
}
