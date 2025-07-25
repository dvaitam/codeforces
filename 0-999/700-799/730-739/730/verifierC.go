package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refC.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "730C.go"))
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(4) + 1
		m := rng.Intn(n*(n-1)/2 + 1)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		edges := make(map[[2]int]struct{})
		for j := 0; j < m; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for u == v || contains(edges, u, v) {
				u = rng.Intn(n) + 1
				v = rng.Intn(n) + 1
			}
			edges[[2]int{u, v}] = struct{}{}
			edges[[2]int{v, u}] = struct{}{}
			fmt.Fprintf(&sb, "%d %d\n", u, v)
		}
		w := rng.Intn(3) + 1
		fmt.Fprintf(&sb, "%d\n", w)
		for j := 0; j < w; j++ {
			c := rng.Intn(n) + 1
			ki := rng.Intn(5) + 1
			p := rng.Intn(10) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", c, ki, p)
		}
		q := rng.Intn(3) + 1
		fmt.Fprintf(&sb, "%d\n", q)
		for j := 0; j < q; j++ {
			g := rng.Intn(n) + 1
			r := rng.Intn(5) + 1
			a := rng.Intn(20) + 1
			fmt.Fprintf(&sb, "%d %d %d\n", g, r, a)
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func contains(m map[[2]int]struct{}, u, v int) bool {
	_, ok := m[[2]int{u, v}]
	return ok
}

func runCase(bin, ref string, c Case) error {
	exp, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(exp) != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
