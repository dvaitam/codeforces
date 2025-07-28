package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1858B.go")
	bin := filepath.Join(os.TempDir(), "oracle1858B.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func deterministicCases() []string {
	return []string{
		"1\n5 2 2\n2 4\n",
		"1\n10 3 3\n1 5 10\n",
		"1\n6 2 2\n3 6\n",
	}
}

func randomCase(rng *rand.Rand) string {
	n := rng.Intn(50) + 2
	m := rng.Intn(5) + 2
	if m > n {
		m = n
	}
	d := rng.Intn(n-1) + 2
	pos := make([]int, m)
	used := make(map[int]bool)
	for i := 0; i < m; i++ {
		for {
			v := rng.Intn(n) + 1
			if !used[v] {
				used[v] = true
				pos[i] = v
				break
			}
		}
	}
	sort.Ints(pos)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n", n, m, d)
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", pos[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, in := range cases {
		want, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := run(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
