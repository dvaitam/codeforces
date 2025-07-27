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
	a, b, c, d int64
	x, y       int64
	x1, y1     int64
	x2, y2     int64
}

func solveCase(tc testCase) string {
	if tc.x1 == tc.x2 && (tc.a > 0 || tc.b > 0) {
		return "NO"
	}
	if tc.y1 == tc.y2 && (tc.c > 0 || tc.d > 0) {
		return "NO"
	}
	nx := tc.x + tc.b - tc.a
	ny := tc.y + tc.d - tc.c
	if nx < tc.x1 || nx > tc.x2 || ny < tc.y1 || ny > tc.y2 {
		return "NO"
	}
	return "YES"
}

func generateCase(rng *rand.Rand) testCase {
	// coordinate range [-10,10]
	x1 := int64(rng.Intn(21) - 10)
	x2 := x1 + int64(rng.Intn(21))
	y1 := int64(rng.Intn(21) - 10)
	y2 := y1 + int64(rng.Intn(21))
	x := x1 + int64(rng.Intn(int(x2-x1+1)))
	y := y1 + int64(rng.Intn(int(y2-y1+1)))
	a := int64(rng.Intn(5))
	b := int64(rng.Intn(5))
	c := int64(rng.Intn(5))
	d := int64(rng.Intn(5))
	return testCase{a, b, c, d, x, y, x1, y1, x2, y2}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("1\n%d %d %d %d\n%d %d %d %d %d %d\n", tc.a, tc.b, tc.c, tc.d, tc.x, tc.y, tc.x1, tc.y1, tc.x2, tc.y2)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solveCase(tc)
	if !strings.EqualFold(got, expected) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			fmt.Fprintf(os.Stderr, "input:\n1\n%d %d %d %d\n%d %d %d %d %d %d\n", tc.a, tc.b, tc.c, tc.d, tc.x, tc.y, tc.x1, tc.y1, tc.x2, tc.y2)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
