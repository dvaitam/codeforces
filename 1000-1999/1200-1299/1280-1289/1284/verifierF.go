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
)

func buildRef() (string, func(), error) {
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1284F.go")
	tmpDir, err := os.MkdirTemp("", "ref1284F")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "refbin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return bin, cleanup, nil
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randTree(r *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := r.Intn(i-1) + 1
		edges = append(edges, [2]int{i, p})
	}
	return edges
}

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 2
	t1 := randTree(r, n)
	t2 := randTree(r, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range t1 {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for _, e := range t2 {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, cleanup, err := buildRef()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()
	r := rand.New(rand.NewSource(1))
	for tc := 1; tc <= 100; tc++ {
		in := genCase(r)
		want, err := run(ref, in)
		if err != nil {
			fmt.Printf("reference failed on case %d: %v\n", tc, err)
			os.Exit(1)
		}
		got, err := run(candidate, in)
		if err != nil {
			fmt.Printf("case %d: %v\n", tc, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %q\ngot: %q\n", tc, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
