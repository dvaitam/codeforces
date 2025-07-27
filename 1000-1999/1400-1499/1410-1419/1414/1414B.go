package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}
func main() {
	in := bufio.NewScanner(os.Stdin)
	in.Split(bufio.ScanWords)
	in.Scan()
	t, _ := strconv.Atoi(in.Text())
	for i := 0; i < t; i++ {
		in.Scan()
		a, _ := strconv.ParseInt(in.Text(), 10, 64)
		in.Scan()
		b, _ := strconv.ParseInt(in.Text(), 10, 64)
		fmt.Println(gcd(a, b))
	}
}
