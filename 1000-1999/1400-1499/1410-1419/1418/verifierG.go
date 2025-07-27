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

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(9) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(5)+1)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func bruteForce(a []int) int {
	n := len(a)
	count := 0
	for l := 0; l < n; l++ {
		freq := make(map[int]int)
		for r := l; r < n; r++ {
			freq[a[r]]++
			if (r-l+1)%3 == 0 {
				ok := true
				for _, v := range freq {
					if v != 3 {
						ok = false
						break
					}
				}
				if ok {
					count++
				}
			}
		}
	}
	return count
}

func expectedOutput(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("bad input")
	}
	var n int
	fmt.Sscan(lines[0], &n)
	nums := strings.Fields(lines[1])
	if len(nums) != n {
		return "", fmt.Errorf("bad input")
	}
	arr := make([]int, n)
	for i, s := range nums {
		fmt.Sscan(s, &arr[i])
	}
	ans := bruteForce(arr)
	return fmt.Sprintf("%d", ans), nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := expectedOutput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error on case %d: %v", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\ngot:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
