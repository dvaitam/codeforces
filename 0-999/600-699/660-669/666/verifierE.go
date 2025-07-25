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
	s       string
	texts   []string
	queries [][4]int
}

func genCase(rng *rand.Rand) testCase {
	slen := rng.Intn(5) + 1
	sb := make([]byte, slen)
	for i := 0; i < slen; i++ {
		sb[i] = byte('a' + rng.Intn(3))
	}
	s := string(sb)
	m := rng.Intn(3) + 1
	texts := make([]string, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(5) + 1
		tb := make([]byte, l)
		for j := 0; j < l; j++ {
			tb[j] = byte('a' + rng.Intn(3))
		}
		texts[i] = string(tb)
	}
	q := rng.Intn(3) + 1
	queries := make([][4]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(m) + 1
		r := l + rng.Intn(m-l+1)
		pl := rng.Intn(slen) + 1
		pr := pl + rng.Intn(slen-pl+1)
		queries[i] = [4]int{l, r, pl, pr}
	}
	return testCase{s: s, texts: texts, queries: queries}
}

func solve(tc testCase) string {
	input := buildInput(tc)
	out, err := run("666E.go", input)
	if err != nil {
		panic(err)
	}
	return out
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(tc.s + "\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.texts)))
	for _, t := range tc.texts {
		sb.WriteString(t + "\n")
	}
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.queries)))
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", q[0], q[1], q[2], q[3]))
	}
	return sb.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		input := buildInput(tc)
		exp := solve(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexp:\n%s\n---\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
