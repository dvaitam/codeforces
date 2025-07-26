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

func expected(arr []int) int {
	n := len(arr)
	b := append([]int(nil), arr...)
	c := 0
	for i := 0; i < n; i++ {
		if b[i] == 0 {
			if i == 0 || i == n-1 {
				continue
			}
			if b[i-1] == 1 && b[i+1] == 1 {
				b[i+1] = 0
				c++
			}
		}
	}
	return c
}

func runCase(exe, input string, exp int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	valStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(valStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", valStr)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tcase := 0; tcase < 100; tcase++ {
		n := rng.Intn(98) + 3
		arr := make([]int, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(2)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(arr)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", tcase+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
