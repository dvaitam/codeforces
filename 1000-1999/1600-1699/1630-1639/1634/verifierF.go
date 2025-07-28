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

func genTests() []string {
	rng := rand.New(rand.NewSource(46))
	tests := make([]string, 100)
	for i := range tests {
		n := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		mod := int64(rng.Intn(1000) + 2)
		var in strings.Builder
		fmt.Fprintf(&in, "%d %d %d\n", n, q, mod)
		for j := 0; j < n; j++ {
			if j > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", rng.Int63n(mod))
		}
		in.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				in.WriteByte(' ')
			}
			fmt.Fprintf(&in, "%d", rng.Int63n(mod))
		}
		in.WriteByte('\n')
		for j := 0; j < q; j++ {
			if rng.Intn(2) == 0 {
				in.WriteString("A ")
			} else {
				in.WriteString("B ")
			}
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&in, "%d %d\n", l, r)
		}
		tests[i] = in.String()
	}
	return tests
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	solPath := "./solF.bin"
	if err := exec.Command("go", "build", "-o", solPath, "1634F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(solPath)

	tests := genTests()
	for i, t := range tests {
		expect, err := run(solPath, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "mismatch on test %d\nexpected:\n%s\nGot:\n%s", i+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
