package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	input string
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1615D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(3)
	tests := make([]Test, 0, 20)
	for t := 0; t < 19; t++ {
		n := rand.Intn(4) + 2 // 2..5
		q := rand.Intn(3)
		edges := make([][3]int, n-1)
		for i := 1; i < n; i++ {
			u := rand.Intn(i) + 1
			v := i + 1
			w := rand.Intn(6) - 1 // maybe -1..4
			if w < 0 {
				w = -1
			}
			edges[i-1] = [3]int{u, v, w}
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
		for i := 0; i < n-1; i++ {
			e := edges[i]
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		for i := 0; i < q; i++ {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			for v == u {
				v = rand.Intn(n) + 1
			}
			w := rand.Intn(2)
			sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
		}
		tests = append(tests, Test{sb.String()})
	}
	// minimal case
	tests = append(tests, Test{"1\n2 0\n1 2 -1\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		want, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
