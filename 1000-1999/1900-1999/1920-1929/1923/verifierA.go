package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	n   int
	arr []int
}

func solve(arr []int) int {
	first := -1
	last := -1
	for i, v := range arr {
		if v == 1 {
			if first == -1 {
				first = i
			}
			last = i
		}
	}
	if first == -1 {
		return 0
	}
	cnt := 0
	for i := first; i <= last; i++ {
		if arr[i] == 0 {
			cnt++
		}
	}
	return cnt
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(49) + 2 // 2..50
		arr := make([]int, n)
		hasOne := false
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 1 {
				arr[i] = 1
				hasOne = true
			}
		}
		if !hasOne {
			arr[rng.Intn(n)] = 1
		}
		tests = append(tests, test{n, arr})
	}
	return tests
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(t.n))
		sb.WriteString("\n")
		for j, v := range t.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteString("\n")
		expected := strconv.Itoa(solve(t.arr))
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
