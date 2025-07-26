package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseA struct {
	n   int
	arr []int
}

func generateCaseA(rng *rand.Rand) (string, testCaseA) {
	n := rng.Intn(100) + 1
	arr := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1000) + 1
		fmt.Fprintf(&sb, "%d ", arr[i])
	}
	sb.WriteByte('\n')
	return sb.String(), testCaseA{n: n, arr: arr}
}

func expectedA(tc testCaseA) string {
	a := tc.arr
	n := len(a)
	i := 0
	for i+1 < n && a[i] < a[i+1] {
		i++
	}
	for i+1 < n && a[i] == a[i+1] {
		i++
	}
	for i+1 < n && a[i] > a[i+1] {
		i++
	}
	if i == n-1 {
		return "YES"
	}
	return "NO"
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseA(rng)
		expect := expectedA(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
