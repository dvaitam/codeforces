package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func sieve(limit int) []int {
	isComp := make([]bool, limit+1)
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					isComp[j] = true
				}
			}
		}
	}
	return primes
}

func countDistinct(n int, primes []int) int {
	cnt := 0
	for _, p := range primes {
		if p*p > n {
			break
		}
		if n%p == 0 {
			cnt++
			for n%p == 0 {
				n /= p
			}
		}
	}
	if n > 1 {
		cnt++
	}
	return cnt
}

func countPairs(c, d, x, g int, primes []int) int {
	y := x / g
	if (y+d)%c != 0 {
		return 0
	}
	k := (y + d) / c
	if k <= 0 {
		return 0
	}
	cnt := countDistinct(k, primes)
	return 1 << cnt
}

func solve(c, d, x int, primes []int) int {
	ans := 0
	for g := 1; g*g <= x; g++ {
		if x%g != 0 {
			continue
		}
		ans += countPairs(c, d, x, g, primes)
		if g*g != x {
			ans += countPairs(c, d, x, x/g, primes)
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	primes := sieve(50000)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		c, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		d, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		x, _ := strconv.Atoi(scan.Text())
		expected[i] = fmt.Sprintf("%d", solve(c, d, x, primes))
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
