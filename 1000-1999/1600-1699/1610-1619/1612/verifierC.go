package main

import (
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func sumFirst(n, k int64) int64 {
	if n <= k {
		return n * (n + 1) / 2
	}
	m := n - k
	p := k - m - 1
	return k*k - p*(p+1)/2
}

func expected(k, x int64) int64 {
	if x >= k*k {
		return 2*k - 1
	}
	low, high := int64(1), 2*k-1
	for low < high {
		mid := (low + high) / 2
		if sumFirst(mid, k) >= x {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func check(k, x int64, output string) error {
	val, err := strconv.ParseInt(strings.TrimSpace(output), 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer output %q", output)
	}
	want := expected(k, x)
	if val != want {
		return fmt.Errorf("expected %d got %d", want, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		k := rand.Int63n(100000) + 1
		maxX := k * k
		x := rand.Int63n(maxX) + 1
		input := fmt.Sprintf("1\n%d %d\n", k, x)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(k, x, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: k=%d x=%d\n", i+1, err, k, x)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
