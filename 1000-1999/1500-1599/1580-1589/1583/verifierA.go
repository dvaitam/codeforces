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

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(8) + 3
	used := make(map[int]bool)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(200) + 1
			if !used[v] {
				used[v] = true
				arr[i] = v
				break
			}
		}
	}
	return arr
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("missing subset size")
	}
	var x int
	if _, err := fmt.Sscan(scanner.Text(), &x); err != nil {
		return fmt.Errorf("bad subset size: %v", err)
	}
	indices := make([]int, 0, x)
	for scanner.Scan() {
		var v int
		if _, err := fmt.Sscan(scanner.Text(), &v); err != nil {
			return fmt.Errorf("bad index: %v", err)
		}
		indices = append(indices, v)
	}
	if len(indices) != x {
		return fmt.Errorf("declared size %d but got %d indices", x, len(indices))
	}
	seen := make(map[int]bool)
	sum := 0
	for _, id := range indices {
		if id < 1 || id > len(arr) {
			return fmt.Errorf("index %d out of range", id)
		}
		if seen[id] {
			return fmt.Errorf("duplicate index %d", id)
		}
		seen[id] = true
		sum += arr[id-1]
	}
	if isPrime(sum) {
		return fmt.Errorf("subset sum %d is prime", sum)
	}
	total := 0
	for _, v := range arr {
		total += v
	}
	expectedSize := len(arr)
	if isPrime(total) {
		expectedSize = len(arr) - 1
	}
	if x != expectedSize {
		return fmt.Errorf("expected subset size %d got %d", expectedSize, x)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		arr := genCase(rng)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
