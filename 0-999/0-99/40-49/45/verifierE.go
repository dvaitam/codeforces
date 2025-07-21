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

func run(bin, input string) (string, error) {
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

func randName(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	b[0] = byte('A' + rng.Intn(26))
	for i := 1; i < l; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	names := make([]string, 0, n)
	surnames := make([]string, 0, n)
	used := make(map[string]bool)
	for len(names) < n {
		s := randName(rng)
		if !used[s] {
			used[s] = true
			names = append(names, s)
		}
	}
	used = make(map[string]bool)
	for len(surnames) < n {
		s := randName(rng)
		if !used[s] {
			used[s] = true
			surnames = append(surnames, s)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, v := range names {
		fmt.Fprintln(&sb, v)
	}
	for _, v := range surnames {
		fmt.Fprintln(&sb, v)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	ref := "45E.go"
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
