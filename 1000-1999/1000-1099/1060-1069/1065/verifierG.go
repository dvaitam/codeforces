package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

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

func fibLengths() [25]int64 {
	var f [25]int64
	f[0] = 1
	f[1] = 1
	for i := 2; i < 25; i++ {
		f[i] = f[i-1] + f[i-2]
	}
	return f
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(48))
	tests := make([]string, 0, 100)
	fl := fibLengths()
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		maxK := fl[n]
		if maxK < 1 {
			maxK = 1
		}
		K := rng.Int63n(maxK) + 1
		m := rng.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, K, m))
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	if cand == "--" && len(os.Args) >= 3 {
		cand = os.Args[2]
	}
	official := "./officialG"
	if err := exec.Command("go", "build", "-o", official, "1065G.go").Run(); err != nil {
		fmt.Println("failed to build official solution:", err)
		os.Exit(1)
	}
	defer os.Remove(official)
	tests := generateTests()
	for i, tc := range tests {
		exp, eerr := run(official, tc)
		got, gerr := run(cand, tc)
		if eerr != nil {
			fmt.Fprintf(os.Stderr, "official failed on test %d: %v\n", i+1, eerr)
			os.Exit(1)
		}
		if gerr != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\ninput:\n%s", i+1, gerr, tc)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
