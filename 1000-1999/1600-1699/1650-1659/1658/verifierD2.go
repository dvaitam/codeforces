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
	ref := "refD2.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1658D2.go")
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
	rng := rand.New(rand.NewSource(5))
	cases := make([]Case, 0, 105)
	for i := 0; i < 100; i++ {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for j := 0; j < t; j++ {
			n := rng.Intn(20) + 1
			l := rng.Intn((1 << 17) - n)
			r := l + n - 1
			x := rng.Intn(1 << 17)
			fmt.Fprintf(&sb, "%d %d\n", l, r)
			for k := 0; k < n; k++ {
				if k > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", (l+k)^x)
			}
			sb.WriteByte('\n')
		}
		cases = append(cases, Case{sb.String()})
	}
	cases = append(cases, Case{"1\n0 3\n7 6 5 4\n"})
	cases = append(cases, Case{"1\n1 4\n4 7 6 5\n"})
	cases = append(cases, Case{"1\n2 4\n3 1 2\n"})
	cases = append(cases, Case{"1\n0 0\n0\n"})
	cases = append(cases, Case{"2\n0 1\n1 0\n0 0\n0\n"})
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
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
