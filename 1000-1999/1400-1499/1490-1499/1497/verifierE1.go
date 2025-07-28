package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func buildPrimes(limit int) []int {
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

func canonical(x int, primes []int) int {
	res := 1
	for _, p := range primes {
		if p*p > x {
			break
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt ^= 1
		}
		if cnt == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func segmentsFor(arr []int, mask int) int {
	seen := map[int]bool{}
	seg := 1
	for i, v := range arr {
		if mask>>i&1 == 1 {
			v = -100000 - i
		}
		if seen[v] {
			seg++
			seen = map[int]bool{}
		}
		seen[v] = true
	}
	return seg
}

func solveBrute(arr []int, k int) int {
	n := len(arr)
	best := n
	for mask := 0; mask < 1<<n; mask++ {
		if bits.OnesCount(uint(mask)) > k {
			continue
		}
		seg := segmentsFor(arr, mask)
		if seg < best {
			best = seg
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE1.txt")
	if err != nil {
		fmt.Println("failed to open testcasesE1.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	primes := buildPrimes(3200)
	scanner := bufio.NewScanner(f)
	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		caseNum++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Printf("case %d invalid format\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		if len(fields)-2 != n {
			fmt.Printf("case %d length mismatch\n", caseNum)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			arr[i] = canonical(v, primes)
		}
		exp := solveBrute(arr, k)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fields[2+i])
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		resultStr := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(resultStr)
		if err != nil {
			fmt.Printf("case %d: invalid output %q\n", caseNum, resultStr)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", caseNum, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}
