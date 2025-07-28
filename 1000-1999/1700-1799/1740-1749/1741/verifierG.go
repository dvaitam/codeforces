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

type Case struct {
	n, m  int
	edges [][2]int
	f     int
	h     []int
	k     int
	p     []int
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 2
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges) + 1
		edgeSet := make(map[[2]int]bool)
		edges := make([][2]int, 0, m)
		for len(edges) < m {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			if a == b {
				continue
			}
			if a > b {
				a, b = b, a
			}
			if edgeSet[[2]int{a, b}] {
				continue
			}
			edgeSet[[2]int{a, b}] = true
			edges = append(edges, [2]int{a, b})
		}
		f := rng.Intn(n) + 1
		h := make([]int, f)
		for j := 0; j < f; j++ {
			h[j] = rng.Intn(n) + 1
		}
		k := rng.Intn(f) + 1
		p := make([]int, k)
		used := rand.Perm(f)[:k]
		for j := 0; j < k; j++ {
			p[j] = used[j] + 1
		}
		cases[i] = Case{n: n, m: m, edges: edges, f: f, h: h, k: k, p: p}
	}
	return cases
}

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refG.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "1741G.go"))
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

func runCase(bin, ref string, c Case) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.m))
	for _, e := range c.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d\n", c.f))
	for i, v := range c.h {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d\n", c.k))
	for i, v := range c.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	input := sb.String()
	exp, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	if exp != got {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %+v\n", i+1, err, c)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
