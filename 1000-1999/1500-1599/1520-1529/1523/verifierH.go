package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsH = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "candH")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func buildOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleH")
	cmd := exec.Command("go", "build", "-o", tmp, "1523H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		maxk := min(30, r-l)
		var k int
		if maxk > 0 {
			k = rng.Intn(maxk + 1)
		}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", l, r, k))
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		return
	}
	bin, clean, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if clean != nil {
		defer clean()
	}
	oracle, c2, err := buildOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c2()

	rng := rand.New(rand.NewSource(7))
	for i := 0; i < numTestsH; i++ {
		input := genCase(rng)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on case %d: %v\n", i+1, err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if want != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
