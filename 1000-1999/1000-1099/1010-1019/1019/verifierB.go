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

func compute(a []int) int {
	n := len(a)
	half := n / 2
	for i := 0; i < half; i++ {
		if a[i] == a[i+half] {
			return i + 1
		}
	}
	return -1
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := 2 * (rng.Intn(5) + 1)
		vals := make([]int, n)
		vals[0] = rng.Intn(10)
		for j := 1; j < n; j++ {
			delta := 1
			if rng.Intn(2) == 0 {
				delta = -1
			}
			vals[j] = vals[j-1] + delta
		}
		// ensure circular difference property
		if abs(vals[0]-vals[n-1]) != 1 {
			if vals[n-1] > vals[0] {
				vals[n-1] = vals[0] - 1
			} else {
				vals[n-1] = vals[0] + 1
			}
		}
		var sb strings.Builder
		fmt.Fprintln(&sb, n)
		for j, v := range vals {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := compute(vals)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		val, err2 := strconv.Atoi(strings.TrimSpace(out))
		if err2 != nil || val != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
