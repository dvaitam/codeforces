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

func isSorted(a []int) bool {
	asc, desc := true, true
	for i := 1; i < len(a); i++ {
		if a[i] < a[i-1] {
			asc = false
		}
		if a[i] > a[i-1] {
			desc = false
		}
		if !asc && !desc {
			return false
		}
	}
	return true
}

func findPair(a []int) (int, int, bool) {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if a[i] == a[j] {
				continue
			}
			b := append([]int(nil), a...)
			b[i], b[j] = b[j], b[i]
			if !isSorted(b) {
				return i + 1, j + 1, true
			}
		}
	}
	return 0, 0, false
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkOutput(arr []int, out string) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if _, _, ok := findPair(arr); ok {
			return fmt.Errorf("solution reported -1 but a valid pair exists")
		}
		return nil
	}
	parts := strings.Fields(out)
	if len(parts) != 2 {
		return fmt.Errorf("expected two integers or -1, got %s", out)
	}
	i1, err1 := strconv.Atoi(parts[0])
	i2, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid integers in output: %s", out)
	}
	n := len(arr)
	if i1 < 1 || i1 > n || i2 < 1 || i2 > n || i1 == i2 {
		return fmt.Errorf("indices out of range or equal")
	}
	if arr[i1-1] == arr[i2-1] {
		return fmt.Errorf("cannot swap equal elements")
	}
	b := append([]int(nil), arr...)
	b[i1-1], b[i2-1] = b[i2-1], b[i1-1]
	if isSorted(b) {
		return fmt.Errorf("array still sorted after swap")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(1000000000)
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprint(v))
		}
		input.WriteByte('\n')
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if err := checkOutput(arr, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
