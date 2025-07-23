package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var k, n, s, p int64
	if _, err := fmt.Fscan(in, &k, &n, &s, &p); err != nil {
		return
	}
	sheetsPerPerson := (n + s - 1) / s
	totalSheets := sheetsPerPerson * k
	packs := (totalSheets + p - 1) / p
	fmt.Println(packs)
}
