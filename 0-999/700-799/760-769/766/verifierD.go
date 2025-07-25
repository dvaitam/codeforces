package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "./refD_bin"
	cmd := exec.Command("go", "build", "-o", ref, "766D.go")
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

func randomWord(rng *rand.Rand, l int) string {
	b := make([]byte, l)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func genTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(5) + 1
		q := rng.Intn(5) + 1
		words := make([]string, n)
		used := map[string]bool{}
		for j := 0; j < n; j++ {
			for {
				w := randomWord(rng, rng.Intn(5)+1)
				if !used[w] {
					used[w] = true
					words[j] = w
					break
				}
			}
		}
		sb := strings.Builder{}
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, q)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(words[j])
		}
		sb.WriteByte('\n')
		for j := 0; j < m; j++ {
			t := rng.Intn(2) + 1
			x := words[rng.Intn(n)]
			y := words[rng.Intn(n)]
			for x == y {
				y = words[rng.Intn(n)]
			}
			fmt.Fprintf(&sb, "%d %s %s\n", t, x, y)
		}
		for j := 0; j < q; j++ {
			x := words[rng.Intn(n)]
			y := words[rng.Intn(n)]
			for x == y {
				y = words[rng.Intn(n)]
			}
			fmt.Fprintf(&sb, "%s %s\n", x, y)
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
