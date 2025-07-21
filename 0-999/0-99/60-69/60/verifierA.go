package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(r *bufio.Reader) string {
	var n, m int
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return ""
	}
	left, right := 1, n
	for i := 0; i < m; i++ {
		var w1, w2, dir, w4 string
		var idx int
		fmt.Fscan(r, &w1, &w2, &dir, &w4, &idx)
		if dir == "left" {
			if idx-1 < right {
				right = idx - 1
			}
		} else {
			if idx+1 > left {
				left = idx + 1
			}
		}
	}
	if left > right {
		return "-1"
	}
	return fmt.Sprintf("%d", right-left+1)
}

func generateCaseA(rng *rand.Rand) string {
	n := rng.Intn(20) + 1 // 1..20
	m := rng.Intn(20)     // 0..19
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		idx := rng.Intn(n) + 1
		if rng.Intn(2) == 0 {
			fmt.Fprintf(&sb, "To the left of %d\n", idx)
		} else {
			fmt.Fprintf(&sb, "To the right of %d\n", idx)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		tc := generateCaseA(rng)
		expect := solveA(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
