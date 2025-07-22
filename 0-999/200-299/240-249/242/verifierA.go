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

type testCase struct {
	x, y, a, b int
}

func solveCase(tc testCase) string {
	var pairs [][2]int
	for c := tc.a; c <= tc.x; c++ {
		for d := tc.b; d <= tc.y && d < c; d++ {
			pairs = append(pairs, [2]int{c, d})
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(pairs))
	for _, p := range pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
	}
	return sb.String()
}

func genCase(rng *rand.Rand) (string, string) {
	tc := testCase{}
	tc.x = rng.Intn(10) + 1
	tc.y = rng.Intn(10) + 1
	tc.a = rng.Intn(tc.x) + 1
	tc.b = rng.Intn(tc.y) + 1
	in := fmt.Sprintf("%d %d %d %d\n", tc.x, tc.y, tc.a, tc.b)
	out := solveCase(tc)
	return in, out
}

func runCase(bin, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got := strings.TrimSpace(buf.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected \n%s\ngot \n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
