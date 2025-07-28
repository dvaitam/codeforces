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
	exe, err := os.CreateTemp("", "refJ-*")
	if err != nil {
		return "", err
	}
	exe.Close()
	path := exe.Name()
	cmd := exec.Command("go", "build", "-o", path, "1725J.go")
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
	for n := 2; n <= 3; n++ {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 2; i <= n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", i-1, i, 1))
		}
		cases = append(cases, Case{sb.String()})
	}
	for len(cases) < 100 {
		n := rng.Intn(5) + 2
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 2; i <= n; i++ {
			u := rng.Intn(i-1) + 1
			w := rng.Int63n(10) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", u, i, w))
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
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
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
