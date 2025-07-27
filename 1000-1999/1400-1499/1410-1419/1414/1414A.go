package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func factorial(n int64) int64 {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	reader.Split(bufio.ScanWords)
	reader.Scan()
	t, _ := strconv.Atoi(reader.Text())
	for i := 0; i < t; i++ {
		reader.Scan()
		nVal, _ := strconv.ParseInt(reader.Text(), 10, 64)
		fmt.Println(factorial(nVal))
	}
}
