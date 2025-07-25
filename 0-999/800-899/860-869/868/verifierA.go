package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	in  string
	out string
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

func solveOracle(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var pass string
	fmt.Fscan(reader, &pass)
	var n int
	fmt.Fscan(reader, &n)
	haveFirst := false
	haveSecond := false
	for i := 0; i < n; i++ {
		var w string
		fmt.Fscan(reader, &w)
		if w == pass {
			return "YES"
		}
		if w[1] == pass[0] {
			haveFirst = true
		}
		if w[0] == pass[1] {
			haveSecond = true
		}
	}
	if haveFirst && haveSecond {
		return "YES"
	}
	return "NO"
}

func randomWord(rng *rand.Rand) string {
	b := []byte{byte('a' + rng.Intn(26)), byte('a' + rng.Intn(26))}
	return string(b)
}

func genCase(rng *rand.Rand) Test {
	pass := randomWord(rng)
	n := rng.Intn(6) + 1
	used := map[string]bool{}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s\n%d\n", pass, n)
	for len(used) < n {
		w := randomWord(rng)
		if !used[w] {
			used[w] = true
			fmt.Fprintln(&sb, w)
		}
	}
	input := sb.String()
	out := solveOracle(input)
	return Test{input, out}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
