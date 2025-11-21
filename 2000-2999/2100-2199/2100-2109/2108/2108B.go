package main

import (
	"bufio"
	"fmt"
	"os"
)

func digitSum(n int) int {
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n%2 == 0 {
		return n == 2
	}
	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func nextPrimeWithDigitSum8(x int) int {
	for {
		if digitSum(x) == 8 && isPrime(x) {
			return x
		}
		x++
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var x int
		fmt.Fscan(in, &x)
		result := nextPrimeWithDigitSum8(x)
		fmt.Fprintln(out, result)
	}
}
