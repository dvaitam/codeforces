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

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func run(bin, input string) (string, error) {
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-n+2) + n - 1 // ensure at least n-1
	if m > maxEdges {
		m = maxEdges
	}
	k := rng.Intn(3) + 1
	edges := make(map[[2]int]struct{})
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)

	// Generate spanning tree to ensure connectivity
	p := rng.Perm(n)
	count := 0
	for i := 0; i < n-1; i++ {
		u, v := p[i], p[i+1]
		if u > v {
			u, v = v, u
		}
		key := [2]int{u + 1, v + 1}
		edges[key] = struct{}{}
		fmt.Fprintf(&sb, "%d %d %d\n", key[0], key[1], rng.Intn(10)+1)
		count++
	}

	for count < m {
		u := rng.Intn(n)
		v := rng.Intn(n)
		if u >= v {
			continue
		}
		key := [2]int{u + 1, v + 1}
		if _, ok := edges[key]; ok {
			continue
		}
		edges[key] = struct{}{}
		fmt.Fprintf(&sb, "%d %d %d\n", key[0], key[1], rng.Intn(10)+1)
		count++
	}
	for i := 0; i < k; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", a, b)
	}
	return sb.String()
}

func runCase(bin, sol, tc string) error {
	want, err := run(sol, tc)
	if err != nil {
		return fmt.Errorf("internal error: %v", err)
	}
	got, err := run(bin, tc)
	if err != nil {
		return err
	}
	if strings.TrimSpace(want) != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(want), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	sol := filepath.Join(baseDir(), "1433G.go")
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, sol, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}