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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(r *bufio.Reader) string {
	var n int
	fmt.Fscan(r, &n)
	var firstSum, secondSum int64
	firstSeq := make([]int64, 0, n)
	secondSeq := make([]int64, 0, n)
	last := 0
	for i := 0; i < n; i++ {
		var ai int64
		fmt.Fscan(r, &ai)
		if ai > 0 {
			firstSum += ai
			firstSeq = append(firstSeq, ai)
			last = 1
		} else {
			val := -ai
			secondSum += val
			secondSeq = append(secondSeq, val)
			last = 2
		}
	}
	if firstSum > secondSum {
		return "first"
	} else if firstSum < secondSum {
		return "second"
	}
	minLen := len(firstSeq)
	if len(secondSeq) < minLen {
		minLen = len(secondSeq)
	}
	for i := 0; i < minLen; i++ {
		if firstSeq[i] > secondSeq[i] {
			return "first"
		} else if firstSeq[i] < secondSeq[i] {
			return "second"
		}
	}
	if len(firstSeq) > len(secondSeq) {
		return "first"
	} else if len(firstSeq) < len(secondSeq) {
		return "second"
	}
	if last == 1 {
		return "first"
	}
	return "second"
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		val := rng.Intn(20) + 1
		if rng.Intn(2) == 0 {
			val = -val
		}
		fmt.Fprintf(&sb, "%d\n", val)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		expect := solveB(bufio.NewReader(strings.NewReader(tc)))
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
