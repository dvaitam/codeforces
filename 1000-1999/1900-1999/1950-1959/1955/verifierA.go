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

func buildRef() (string, error) {
	ref := "refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1955A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v: %s", err, string(out))
	}
	return "./" + ref, nil
}

func runProgram(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(100) + 1
	a := rng.Intn(30) + 1
	b := rng.Intn(30) + 1
	return fmt.Sprintf("1\n%d %d %d\n", n, a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 0, 100)
	// Some edge cases
	cases = append(cases, "1\n1 1 1\n", "1\n2 5 3\n", "1\n100 30 30\n")
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		exp, err := runProgram(ref, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProgram(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		exp = strings.TrimSpace(exp)
		got = strings.TrimSpace(got)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
