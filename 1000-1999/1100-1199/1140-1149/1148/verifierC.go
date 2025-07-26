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

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase() (string, []int) {
	n := (rand.Intn(4) + 1) * 2 // even 2..8
	perm := rand.Perm(n)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = perm[i] + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), arr
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func check(n int, arr []int, out string) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid first line")
	}
	if m < 0 || m > 5*n {
		return fmt.Errorf("invalid operations count")
	}
	if len(lines)-1 != m {
		return fmt.Errorf("expected %d operation lines", m)
	}
	for i := 1; i <= m; i++ {
		parts := strings.Fields(lines[i])
		if len(parts) != 2 {
			return fmt.Errorf("op %d: need two numbers", i)
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("op %d: index out of range", i)
		}
		if 2*abs(a-b) < n {
			return fmt.Errorf("op %d: distance too small", i)
		}
		arr[a-1], arr[b-1] = arr[b-1], arr[a-1]
	}
	for i := 0; i < n; i++ {
		if arr[i] != i+1 {
			return fmt.Errorf("array not sorted")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		input, arr := genCase()
		arrCopy := make([]int, len(arr))
		copy(arrCopy, arr)
		out, err := run(bin, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if err := check(len(arr), arrCopy, out); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%soutput:\n%s", t+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
