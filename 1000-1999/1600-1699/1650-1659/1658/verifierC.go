package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Case struct{ input string }

func buildRef() (string, error) {
	ref := "refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1658C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runExe(bin, input string) (string, error) {
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

func genCases() []Case {
	rng := rand.New(rand.NewSource(3))
	cases := make([]Case, 0, 105)
	for i := 0; i < 100; i++ {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for j := 0; j < t; j++ {
			n := rng.Intn(20) + 1
			fmt.Fprintf(&sb, "%d\n", n)
			for k := 0; k < n; k++ {
				if k > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", rng.Intn(n)+1)
			}
			sb.WriteByte('\n')
		}
		cases = append(cases, Case{sb.String()})
	}
	cases = append(cases, Case{"1\n1\n1\n"})
	cases = append(cases, Case{"1\n2\n1 1\n"})
	cases = append(cases, Case{"1\n3\n1 2 3\n"})
	cases = append(cases, Case{"1\n6\n2 3 1 2 3 4\n"})
	cases = append(cases, Case{"2\n1\n1\n1\n1\n"})
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
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
		exp, err := runExe(ref, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, c.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
