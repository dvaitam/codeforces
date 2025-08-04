package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// smallestPrime returns the smallest prime factor of x.
func smallestPrime(x int) int {
	for p := 2; p*p <= x; p++ {
		if x%p == 0 {
			return p
		}
	}
	return x
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("invalid test %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+n {
			fmt.Printf("test %d wrong count\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+i])
			arr[i] = v
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		parts := strings.Split(outStr, "\n")
		if len(parts) != 2 {
			fmt.Printf("Test %d: expected two lines got %d\n", idx, len(parts))
			os.Exit(1)
		}
		gotA := strings.Fields(parts[0])
		gotB := strings.Fields(parts[1])
		if len(gotA) != n || len(gotB) != n {
			fmt.Printf("Test %d: wrong output length\n", idx)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			ga, _ := strconv.Atoi(gotA[i])
			gb, _ := strconv.Atoi(gotB[i])
			p := smallestPrime(arr[i])
			x := arr[i]
			for x%p == 0 {
				x /= p
			}
			if x == 1 {
				if ga != -1 || gb != -1 {
					fmt.Printf("Test %d failed at index %d expected -1 -1 got %d %d\n", idx, i, ga, gb)
					os.Exit(1)
				}
			} else {
				if ga <= 1 || gb <= 1 {
					fmt.Printf("Test %d failed at index %d: divisors must be >1\n", idx, i)
					os.Exit(1)
				}
				if arr[i]%(ga*gb) != 0 {
					fmt.Printf("Test %d failed at index %d: product %d*%d does not divide %d\n", idx, i, ga, gb, arr[i])
					os.Exit(1)
				}
				if gcd(ga, gb) != 1 {
					fmt.Printf("Test %d failed at index %d: gcd(%d,%d) != 1\n", idx, i, ga, gb)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
