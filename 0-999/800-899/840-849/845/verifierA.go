package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Case struct {
	input string
}

func buildRef() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "845A.go")
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(8451))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(100) + 1
		arr := make([]int, 2*n)
		for j := 0; j < 2*n; j++ {
			arr[j] = rng.Intn(1000) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	expected, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(expected) != strings.TrimSpace(got) {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
