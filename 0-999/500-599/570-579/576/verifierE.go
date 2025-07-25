package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "./refE_bin"
	cmd := exec.Command("go", "build", "-o", ref, "576E.go")
	cmd.Stdout = new(bytes.Buffer)
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return ref, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	err := cmd.Run()
	return out.String(), err
}

func genTests() []string {
	rand.Seed(5)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 2
		m := rand.Intn(10) + 1
		k := rand.Intn(3) + 1
		q := rand.Intn(10) + 1
		sb := strings.Builder{}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, k, q))
		type edge struct{ a, b int }
		used := make(map[edge]bool)
		for j := 0; j < m; j++ {
			var a, b int
			for {
				a = rand.Intn(n) + 1
				b = rand.Intn(n) + 1
				if a != b && !used[edge{a, b}] && !used[edge{b, a}] {
					used[edge{a, b}] = true
					break
				}
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
		}
		for j := 0; j < q; j++ {
			e := rand.Intn(m) + 1
			c := rand.Intn(k) + 1
			sb.WriteString(fmt.Sprintf("%d %d\n", e, c))
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for idx, input := range tests {
		exp, err1 := runBinary(ref, input)
		out, err2 := runBinary(cand, input)
		if err1 != nil || err2 != nil {
			fmt.Printf("runtime error on test %d\n", idx+1)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(out) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s", idx+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
