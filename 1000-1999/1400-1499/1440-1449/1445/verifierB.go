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

type testCaseB struct {
	a int
	b int
	c int
	d int
}

func generateB(rng *rand.Rand) testCaseB {
	a := rng.Intn(10)    // 0..9
	d := rng.Intn(a + 1) // 0..a
	c := rng.Intn(10)    // 0..9
	if c < d {           // ensure b <= c after generation
		c = d
	}
	b := rng.Intn(c-d+1) + d // b in [d, c]
	return testCaseB{a: a, b: b, c: c, d: d}
}

func expectedB(tc testCaseB) string {
	sum1 := tc.a + tc.b
	sum2 := tc.c + tc.d
	if sum2 > sum1 {
		sum1 = sum2
	}
	return fmt.Sprintf("%d", sum1)
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, tc testCaseB) error {
	input := fmt.Sprintf("1\n%d %d %d %d\n", tc.a, tc.b, tc.c, tc.d)
	expected := expectedB(tc)
	got, err := run(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateB(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
