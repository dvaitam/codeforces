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

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1611E1.go")
	bin := filepath.Join(os.TempDir(), "oracle1611E1.bin")
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genTree(r *rand.Rand, n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := r.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return edges
}

func genCase(r *rand.Rand) string {
	t := r.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(6) + 1
		k := r.Intn(n) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		used := make(map[int]bool)
		for j := 0; j < k; j++ {
			x := r.Intn(n) + 1
			for used[x] {
				x = r.Intn(n) + 1
			}
			used[x] = true
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", x)
		}
		sb.WriteByte('\n')
		edges := genTree(r, n)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
