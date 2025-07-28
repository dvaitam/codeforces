package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
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
	if errBuf.Len() > 0 {
		return "", fmt.Errorf(errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(arr []int) int {
	sort.Ints(arr)
	n := len(arr)
	return 2*(arr[n-1]-arr[0]) + 2*(arr[n-2]-arr[1])
}

func genTests() [][]int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([][]int, 0, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(97) + 4 // 4..100
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(2000001) - 1000000
		}
		tests = append(tests, a)
	}
	// edge cases
	tests = append(tests, []int{0, 0, 0, 0})
	tests = append(tests, []int{1, -1, 1, -1})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for idx, arr := range tests {
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
		expect := fmt.Sprintf("%d", expected(append([]int(nil), arr...)))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
