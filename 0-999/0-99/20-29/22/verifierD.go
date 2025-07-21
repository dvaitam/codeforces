package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type segment struct{ l, r int }

func solveCase(segs []segment) string {
	sort.Slice(segs, func(i, j int) bool { return segs[i].r < segs[j].r })
	nails := make([]int, 0, len(segs))
	hasLast := false
	last := 0
	for _, s := range segs {
		if !hasLast || last < s.l {
			last = s.r
			nails = append(nails, last)
			hasLast = true
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(nails)))
	for i, v := range nails {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return strings.TrimSpace(sb.String())
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(6) + 1
		segs := make([]segment, n)
		for i := 0; i < n; i++ {
			x := rng.Intn(21) - 10
			y := rng.Intn(21) - 10
			if x <= y {
				segs[i] = segment{x, y}
			} else {
				segs[i] = segment{y, x}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, s := range segs {
			sb.WriteString(fmt.Sprintf("%d %d\n", s.l, s.r))
		}
		input := sb.String()
		expected := solveCase(segs)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
