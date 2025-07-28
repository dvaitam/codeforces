package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const numTestsF = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), fmt.Sprintf("candF_%d", time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("compile candidate: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, func() {}, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), fmt.Sprintf("oracleF_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", tmp, "1675F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("compile oracle: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 2
	k := rng.Intn(n) + 1
	x := rng.Intn(n) + 1
	y := rng.Intn(n) + 1
	tasks := make([]int, k)
	for i := 0; i < k; i++ {
		tasks[i] = rng.Intn(n) + 1
	}
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{parent, i})
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	for i := 0; i < k; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", tasks[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	cand, cleanCand, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanCand()

	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanOracle()

	rng := rand.New(rand.NewSource(6))
	for i := 0; i < numTestsF; i++ {
		input := genCase(rng)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
