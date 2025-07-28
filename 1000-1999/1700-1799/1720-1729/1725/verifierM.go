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
	exe, err := os.CreateTemp("", "refM-*")
	if err != nil {
		return "", err
	}
	exe.Close()
	path := exe.Name()
	cmd := exec.Command("go", "build", "-o", path, "1725M.go")
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
	for _, n := range []int{2, 3} {
		var sb strings.Builder
		m := n - 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 2; i <= n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", i-1, i, 1))
		}
		cases = append(cases, Case{sb.String()})
	}
	for len(cases) < 100 {
		n := rng.Intn(4) + 2
		m := rng.Intn(n*(n-1)) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		edges := make(map[[2]int]bool)
		for i := 0; i < m; i++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n-1) + 1
			if v >= u {
				v++
			}
			key := [2]int{u, v}
			if edges[key] {
				i--
				continue
			}
			edges[key] = true
			w := rng.Int63n(10) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
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
		fmt.Println("usage: go run verifierM.go /path/to/binary")
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
