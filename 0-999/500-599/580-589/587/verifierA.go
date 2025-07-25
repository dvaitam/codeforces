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

func expected(weights []int) int {
	maxW := 0
	for _, w := range weights {
		if w > maxW {
			maxW = w
		}
	}
	cnt := make([]int, maxW+70)
	for _, w := range weights {
		cnt[w]++
	}
	for i := 0; i < len(cnt)-1; i++ {
		pairs := cnt[i] / 2
		if pairs > 0 {
			cnt[i+1] += pairs
			cnt[i] %= 2
		}
	}
	steps := 0
	for _, c := range cnt {
		steps += c
	}
	return steps
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
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(20) + 1
		weights := make([]int, n)
		for i := 0; i < n; i++ {
			weights[i] = rng.Intn(60)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, w := range weights {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", w))
		}
		sb.WriteString("\n")
		input := sb.String()
		expectedStr := fmt.Sprintf("%d", expected(weights))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != expectedStr {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, expectedStr, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
