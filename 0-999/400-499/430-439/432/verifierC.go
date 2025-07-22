package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func sieve(n int) []bool {
	prime := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		prime[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if prime[i] {
			for j := i * i; j <= n; j += i {
				prime[j] = false
			}
		}
	}
	return prime
}

func runCase(exe string, arr []int, primes []bool) error {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if m > 5*n {
		return fmt.Errorf("too many operations: %d", m)
	}
	ops := make([][2]int, m)
	for i := 0; i < m; i++ {
		if _, err := fmt.Fscan(reader, &ops[i][0], &ops[i][1]); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if ops[i][0] < 1 || ops[i][1] < ops[i][0] || ops[i][1] > n {
			return fmt.Errorf("invalid indices in operation")
		}
		if ops[i][1]-ops[i][0]+1 >= len(primes) || !primes[ops[i][1]-ops[i][0]+1] {
			return fmt.Errorf("swap length not prime")
		}
	}
	arrCopy := append([]int(nil), arr...)
	for _, op := range ops {
		i, j := op[0]-1, op[1]-1
		arrCopy[i], arrCopy[j] = arrCopy[j], arrCopy[i]
	}
	for i := 0; i < n; i++ {
		if arrCopy[i] != i+1 {
			return fmt.Errorf("array not sorted")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	primes := sieve(1000)
	for i := 0; i < 100; i++ {
		n := rng.Intn(30) + 1
		arr := rand.Perm(n)
		for j := 0; j < n; j++ {
			arr[j]++
		}
		if err := runCase(exe, arr, primes); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
