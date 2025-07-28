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

type result struct {
	out string
	err error
}

func buildRef() (string, error) {
	refSrc := filepath.Join(filepath.Dir(os.Args[0]), "1901E.go")
	bin := filepath.Join(os.TempDir(), fmt.Sprintf("refE_%d", time.Now().UnixNano()))
	cmd := exec.Command("go", "build", "-o", bin, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return bin, nil
}

func runBinary(path, input string) result {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, stderr.String())
	}
	return result{strings.TrimSpace(out.String()), err}
}

func genTreeEdges(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func genTest() string {
	n := rand.Intn(8) + 2
	vals := make([]int64, n)
	for i := 0; i < n; i++ {
		vals[i] = rand.Int63n(2001) - 1000
	}
	edges := genTreeEdges(n)
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "1\n%d\n", n)
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(5)
	refBin, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)
	for i := 0; i < 100; i++ {
		tc := genTest()
		exp := runBinary(refBin, tc)
		if exp.err != nil {
			fmt.Fprintf(os.Stderr, "reference run failed on test %d: %v\n", i+1, exp.err)
			os.Exit(1)
		}
		got := runBinary(binary, tc)
		if got.err != nil {
			fmt.Fprintf(os.Stderr, "binary failed on test %d: %v\n", i+1, got.err)
			os.Exit(1)
		}
		if exp.out != got.out {
			fmt.Printf("mismatch on test %d\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc, exp.out, got.out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
