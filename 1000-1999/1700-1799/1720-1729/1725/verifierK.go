package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	exe, err := os.CreateTemp("", "refK-*")
	if err != nil {
		return "", err
	}
	exe.Close()
	path := exe.Name()
	cmd := exec.Command("go", "build", "-o", path, "1725K.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
}

func runProg(bin, input string) (string, error) {
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

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(1725))
	cases := make([]Case, 0, 100)
	for _, n := range []int{1, 2} {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d ", rng.Intn(5)+1))
		}
		sb.WriteByte('\n')
		sb.WriteString("1\n2 1\n")
		cases = append(cases, Case{sb.String()})
	}
	for len(cases) < 100 {
		N := rng.Intn(4) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", N))
		for i := 0; i < N; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(10)+1))
		}
		sb.WriteByte('\n')
		Q := rng.Intn(3) + 1
		sb.WriteString(fmt.Sprintf("%d\n", Q))
		for i := 0; i < Q; i++ {
			typ := rng.Intn(3) + 1
			if typ == 1 {
				k := rng.Intn(N) + 1
				w := rng.Intn(10) + 1
				sb.WriteString(fmt.Sprintf("1 %d %d\n", k, w))
			} else if typ == 2 {
				k := rng.Intn(N) + 1
				sb.WriteString(fmt.Sprintf("2 %d\n", k))
			} else {
				l := rng.Intn(N) + 1
				r := rng.Intn(N-l+1) + l
				sb.WriteString(fmt.Sprintf("3 %d %d\n", l, r))
			}
		}
		cases = append(cases, Case{sb.String()})
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	expect, err := runProg(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runProg(bin, c.input)
	if err != nil {
		return err
	}
	if expect != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
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
