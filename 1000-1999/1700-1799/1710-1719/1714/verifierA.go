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

type alarm struct{ h, m int }

func run(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func solve(n, H, M int, al []alarm) (int, int) {
	start := H*60 + M
	best := 24 * 60
	for _, a := range al {
		cur := a.h*60 + a.m
		diff := cur - start
		if diff < 0 {
			diff += 24 * 60
		}
		if diff < best {
			best = diff
		}
	}
	return best / 60, best % 60
}

func generateCase(rng *rand.Rand) (string, int, int, []alarm) {
	n := rng.Intn(10) + 1
	H := rng.Intn(24)
	M := rng.Intn(60)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, H, M))
	alarms := make([]alarm, n)
	for i := 0; i < n; i++ {
		h := rng.Intn(24)
		m := rng.Intn(60)
		sb.WriteString(fmt.Sprintf("%d %d\n", h, m))
		alarms[i] = alarm{h, m}
	}
	return sb.String(), H, M, alarms
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, H, M, alarms := generateCase(rng)
		expH, expM := solve(len(alarms), H, M, alarms)
		expected := fmt.Sprintf("%d %d", expH, expM)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
