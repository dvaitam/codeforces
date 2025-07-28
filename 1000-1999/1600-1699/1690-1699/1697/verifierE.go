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
	input string
}

func buildRef() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1697E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
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

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	coords := make(map[[2]int]struct{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for len(coords) < n {
		x := rng.Intn(11) - 5
		y := rng.Intn(11) - 5
		key := [2]int{x, y}
		if _, ok := coords[key]; ok {
			continue
		}
		coords[key] = struct{}{}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return testCase{input: sb.String()}
}

func runCase(bin, ref string, tc testCase) error {
	expected, err := runBinary(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, tc.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
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

	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		cases[i] = genCase(rng)
	}

	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
