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
	src := filepath.Join(dir, "1917D.go")
	bin := filepath.Join(os.TempDir(), "oracle1917D.bin")
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
		"1\n1 1\n1\n0\n",
		"1\n2 2\n1 3\n0 1\n",
		"1\n3 2\n3 5 1\n1 0\n",
	}
}

func randPerm(rng *rand.Rand, n int) []int {
	a := make([]int, n)
	for i := range a {
		a[i] = i
	}
	rng.Shuffle(n, func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}

func randOddPerm(rng *rand.Rand, n int) []int {
	a := make([]int, n)
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = 2*i + 1
	}
	rng.Shuffle(n, func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	copy(a, vals)
	return a
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	k := rng.Intn(4) + 1
	p := randOddPerm(rng, n)
	q := randPerm(rng, k)
	b := strings.Builder{}
	b.WriteString("1\n")
	b.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, x := range p {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprint(x))
	}
	b.WriteString("\n")
	for i, x := range q {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprint(x))
	}
	b.WriteString("\n")
	return b.String()
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

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, in := range cases {
		want, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := run(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\n got: %s\n", i+1, in, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
