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
	exe, err := os.CreateTemp("", "refB-*")
	if err != nil {
		return "", err
	}
	exe.Close()
	path := exe.Name()
	cmd := exec.Command("go", "build", "-o", path, "1725B.go")
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
	preset := []struct {
		n   int
		d   int64
		arr []int64
	}{
		{1, 1, []int64{1}},
		{2, 1, []int64{1, 2}},
		{3, 10, []int64{5, 5, 5}},
	}
	for _, p := range preset {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", p.n, p.d))
		for i, v := range p.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		cases = append(cases, Case{sb.String()})
	}
	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		d := rng.Int63n(1_000_000_000) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			v := rng.Int63n(1_000_000_000) + 1
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		cases = append(cases, Case{sb.String()})
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	expected, err := runProg(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runProg(bin, c.input)
	if err != nil {
		return err
	}
	if expected != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
