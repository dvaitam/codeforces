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
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "788E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return "./" + ref, nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	input string
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	skills := make([]int, n)
	for i := range skills {
		skills[i] = rng.Intn(20) + 1
	}
	q := rng.Intn(8) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range skills {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		typ := rng.Intn(2) + 1
		x := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", typ, x)
	}
	return testCase{input: sb.String()}
}

func runCase(candidate, reference string, tc testCase) error {
	exp, err := run(reference, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	out, err := run(candidate, tc.input)
	if err != nil {
		return err
	}
	if out != exp {
		return fmt.Errorf("expected %s got %s\ninput:\n%s", exp, out, tc.input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
	var cases []testCase
	cases = append(cases, genCase(rng))
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
