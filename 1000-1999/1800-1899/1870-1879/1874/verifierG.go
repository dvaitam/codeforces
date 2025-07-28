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
	_, cur, _, _ := runtime.Caller(0)
	dir := filepath.Dir(cur)
	src := filepath.Join(dir, "1874G.go")
	bin := filepath.Join(os.TempDir(), "1874G_ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runProg(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) []byte {
	n := rng.Intn(5) + 2 // at least 2
	// base edges to ensure path
	edges := make(map[[2]int]struct{})
	for i := 1; i < n; i++ {
		edges[[2]int{i, i + 1}] = struct{}{}
	}
	extra := rng.Intn(3)
	for k := 0; k < extra; k++ {
		u := rng.Intn(n-1) + 1
		v := rng.Intn(n-u) + u + 1
		edges[[2]int{u, v}] = struct{}{}
	}
	m := len(edges)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 1; i <= n; i++ {
		if i == 1 || i == n {
			sb.WriteString("0\n")
			continue
		}
		typ := rng.Intn(5)
		switch typ {
		case 1:
			a := rng.Intn(5) + 1
			b := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d\n", a, b))
		case 2:
			x := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("2 %d\n", x))
		case 3:
			y := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("3 %d\n", y))
		case 4:
			w := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("4 %d\n", w))
		default:
			sb.WriteString("0\n")
		}
	}
	for e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		want, err := runProg(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runProg(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s", i, err, string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, string(input), want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
