package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func score(m, w [5]int, hs, hu int) int {
	x := [5]int{500, 1000, 1500, 2000, 2500}
	total := 0
	for i := 0; i < 5; i++ {
		score1 := 3 * x[i] / 10
		score2 := (250-m[i])*x[i]/250 - 50*w[i]
		if score1 > score2 {
			total += score1
		} else {
			total += score2
		}
	}
	total += hs * 100
	total -= hu * 50
	return total
}

func generateCase(rng *rand.Rand) (string, string) {
	var m, w [5]int
	for i := 0; i < 5; i++ {
		m[i] = rng.Intn(120)
	}
	for i := 0; i < 5; i++ {
		w[i] = rng.Intn(11)
	}
	hs := rng.Intn(21)
	hu := rng.Intn(21)
	input := fmt.Sprintf("%d %d %d %d %d\n%d %d %d %d %d\n%d %d\n", m[0], m[1], m[2], m[3], m[4], w[0], w[1], w[2], w[3], w[4], hs, hu)
	exp := fmt.Sprintf("%d", score(m, w, hs, hu))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fixed := []struct {
		m, w   [5]int
		hs, hu int
	}{
		{[5]int{0, 0, 0, 0, 0}, [5]int{0, 0, 0, 0, 0}, 0, 0},
		{[5]int{119, 119, 119, 119, 119}, [5]int{10, 10, 10, 10, 10}, 20, 20},
		{[5]int{30, 50, 70, 90, 110}, [5]int{1, 2, 3, 4, 5}, 5, 3},
	}

	caseNum := 1
	for _, tc := range fixed {
		input := fmt.Sprintf("%d %d %d %d %d\n%d %d %d %d %d\n%d %d\n", tc.m[0], tc.m[1], tc.m[2], tc.m[3], tc.m[4], tc.w[0], tc.w[1], tc.w[2], tc.w[3], tc.w[4], tc.hs, tc.hu)
		exp := fmt.Sprintf("%d", score(tc.m, tc.w, tc.hs, tc.hu))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum, exp, out, input)
			os.Exit(1)
		}
		caseNum++
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum, exp, out, in)
			os.Exit(1)
		}
		caseNum++
	}
	fmt.Println("All tests passed")
}
