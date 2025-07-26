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

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1081D.go")
	bin := filepath.Join(os.TempDir(), "oracle1081D.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	n := r.Intn(6) + 2
	m := n - 1 + r.Intn(n)
	k := r.Intn(n-1) + 1
	// ensure k unique specials
	specials := make([]int, k)
	used := make(map[int]bool)
	for i := 0; i < k; i++ {
		for {
			v := r.Intn(n) + 1
			if !used[v] {
				used[v] = true
				specials[i] = v
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < k; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", specials[i])
	}
	sb.WriteByte('\n')
	edges := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			i--
			continue
		}
		if u > v {
			u, v = v, u
		}
		if edges[[2]int{u, v}] {
			i--
			continue
		}
		edges[[2]int{u, v}] = true
		w := r.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", u, v, w)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	r := rand.New(rand.NewSource(1))
	cases := []string{
		"2 1 1\n1\n1 2 5\n",
		"3 3 2\n1 2\n1 2 1\n2 3 2\n1 3 3\n",
		"4 4 2\n1 3\n1 2 1\n2 3 2\n3 4 3\n1 4 4\n",
	}
	for i := 0; i < 97; i++ {
		cases = append(cases, genCase(r))
	}
	for idx, input := range cases {
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
