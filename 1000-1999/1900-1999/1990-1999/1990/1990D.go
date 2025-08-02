package main

import (
	"bufio"
	"fmt"
	"os"
)

const N = 200005

func solve(scanner *bufio.Scanner, writer *bufio.Writer) {
	scanner.Scan()
	var n int
	fmt.Sscan(scanner.Text(), &n)
	
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		scanner.Scan()
		fmt.Sscan(scanner.Text(), &a[i])
		if a[i]%2 == 1 {
			a[i]++
		}
	}
	
	t := 0
	for i := 1; i <= n; i++ {
		if a[i] == 4 && a[i-1] == 4 {
			a[i], a[i-1] = 2, 2
			t++
		}
	}
	
	for i := 1; i <= n; i++ {
		if a[i] == 2 && a[i-1] == 2 {
			a[i], a[i-1] = 0, 0
			t++
		}
	}
	
	for i := 1; i <= n; i++ {
		if a[i] != 0 {
			t++
		}
	}
	
	fmt.Fprintln(writer, t)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	
	scanner.Split(bufio.ScanWords)
	
	scanner.Scan()
	var T int
	fmt.Sscan(scanner.Text(), &T)
	
	for T > 0 {
		solve(scanner, writer)
		T--
	}
}
