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
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1187B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildName(rng *rand.Rand, s string) string {
	counts := make(map[byte]int)
	for i := 0; i < len(s); i++ {
		counts[s[i]]++
	}
	length := rng.Intn(len(s)) + 1
	res := make([]byte, length)
	letters := make([]byte, 0, len(counts))
	for c := range counts {
		letters = append(letters, c)
	}
	for i := 0; i < length; i++ {
		for {
			c := letters[rng.Intn(len(letters))]
			if counts[c] > 0 {
				counts[c]--
				res[i] = c
				break
			}
		}
	}
	return string(res)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d\n", n))
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = letters[rng.Intn(len(letters))]
	}
	sb.WriteString(string(s) + "\n")
	m := rng.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		name := buildName(rng, string(s))
		sb.WriteString(name + "\n")
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		out, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
