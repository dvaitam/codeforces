package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) (int, int) {
	n := rng.Intn(49) + 2 // 2..50
	k := rng.Intn(n-1) + 1
	return n, k
}

func check(n, k int, out string) error {
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanWords)
	arr := make([]int, 0, n)
	for scan.Scan() {
		if len(arr) == n {
			return fmt.Errorf("extra output")
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return fmt.Errorf("invalid integer")
		}
		arr = append(arr, v)
	}
	if len(arr) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(arr))
	}
	seen := make([]bool, n+1)
	for _, v := range arr {
		if v < 1 || v > n {
			return fmt.Errorf("value out of range")
		}
		if seen[v] {
			return fmt.Errorf("duplicate value")
		}
		seen[v] = true
	}
	diffs := make(map[int]struct{})
	for i := 0; i < n-1; i++ {
		d := arr[i] - arr[i+1]
		if d < 0 {
			d = -d
		}
		diffs[d] = struct{}{}
	}
	if len(diffs) != k {
		return fmt.Errorf("expected %d distinct diffs got %d", k, len(diffs))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k := generateCase(rng)
		input := fmt.Sprintf("%d %d\n", n, k)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(n, k, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
