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

func solveA(a [5]int, s string) string {
	total := 0
	for _, ch := range s {
		idx := int(ch - '0')
		if idx >= 1 && idx <= 4 {
			total += a[idx]
		}
	}
	return fmt.Sprintf("%d", total)
}

func generateCase(rng *rand.Rand) (string, string) {
	var a [5]int
	for i := 1; i <= 4; i++ {
		a[i] = rng.Intn(1000)
	}
	l := rng.Intn(20) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('1' + rng.Intn(4))
	}
	s := string(b)
	input := fmt.Sprintf("%d %d %d %d\n%s\n", a[1], a[2], a[3], a[4], s)
	return input, solveA(a, s)
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
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

	cases := []struct{ in, out string }{
		{"1 1 1 1\n1111\n", "4"},
		{"10 20 30 40\n2143\n", "100"},
	}
	for i := 0; i < 100; i++ {
		in, out := generateCase(rng)
		cases = append(cases, struct{ in, out string }{in, out})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.in, tc.out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
