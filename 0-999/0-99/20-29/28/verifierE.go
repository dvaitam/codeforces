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

func referenceSolve(input string) (string, error) {
	cmd := exec.Command("go", "run", "28E.go")
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func verifyE(input, output string) error {
	expect, err := referenceSolve(input)
	if err != nil {
		return fmt.Errorf("reference solution failed: %v", err)
	}
	if strings.TrimSpace(expect) != strings.TrimSpace(output) {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(expect), strings.TrimSpace(output))
	}
	return nil
}

func runCase(bin, tc string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verifyE(tc, out.String())
}

func generateCase(rng *rand.Rand) string {
	// simple triangle polygon
	n := 3
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	coords := [][2]int{{0, 0}, {10, 0}, {0, 10}}
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d %d\n", coords[i][0], coords[i][1])
	}
	fmt.Fprintf(&b, "0 0\n")
	fmt.Fprintf(&b, "1 0 1\n")
	fmt.Fprintf(&b, "1\n")
	fmt.Fprintf(&b, "0 1 -1\n")
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
