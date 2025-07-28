package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1635F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func buildCase(n, q int, xs, ws []int64, qs [][2]int) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ws[i]))
	}
	for _, qu := range qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	return []byte(sb.String())
}

func genRandomCase(rng *rand.Rand) []byte {
	n := rng.Intn(5) + 2
	q := rng.Intn(5) + 1
	xs := make([]int64, n)
	ws := make([]int64, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(10) + 1)
		xs[i] = cur
		ws[i] = int64(rng.Intn(10) + 1)
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		l := rng.Intn(n-1) + 1
		r := rng.Intn(n-l) + l + 1
		queries[i] = [2]int{l, r}
	}
	return buildCase(n, q, xs, ws, queries)
}

func genTests() [][]byte {
	rng := rand.New(rand.NewSource(47))
	tests := make([][]byte, 0, 100)
	xs := []int64{0, 1}
	ws := []int64{1, 1}
	tests = append(tests, buildCase(2, 1, xs, ws, [][2]int{{1, 2}}))
	for len(tests) < 100 {
		tests = append(tests, genRandomCase(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n%s", i+1, err, exp)
			os.Exit(1)
		}
		got, err := runExe(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, string(tc), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
