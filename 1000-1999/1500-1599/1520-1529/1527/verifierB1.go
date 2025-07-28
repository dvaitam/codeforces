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

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1527B1.go")
	bin := filepath.Join(os.TempDir(), "ref1527B1.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return bin, nil
}

func runBinary(bin string, input []byte) ([]byte, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.Bytes(), nil
}

type Case struct{ input []byte }

func genPalin(n int, rng *rand.Rand) string {
	b := make([]byte, n)
	for i := 0; i < (n+1)/2; i++ {
		bit := byte('0' + rng.Intn(2))
		b[i] = bit
		b[n-1-i] = bit
	}
	hasZero := false
	for i := 0; i < n; i++ {
		if b[i] == '0' {
			hasZero = true
			break
		}
	}
	if !hasZero {
		pos := rng.Intn(n)
		b[pos] = '0'
		b[n-1-pos] = '0'
	}
	return string(b)
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 0, 101)
	// simple deterministic case
	{
		var sb strings.Builder
		sb.WriteString("1\n1\n0\n")
		cases = append(cases, Case{[]byte(sb.String())})
	}
	for i := 0; i < 100; i++ {
		t := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for j := 0; j < t; j++ {
			n := rng.Intn(10) + 1
			s := genPalin(n, rng)
			fmt.Fprintf(&sb, "%d\n%s\n", n, s)
		}
		cases = append(cases, Case{[]byte(sb.String())})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, c := range cases {
		exp, err := runBinary(ref, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(string(out)) != strings.TrimSpace(string(exp)) {
			fmt.Printf("wrong answer on case %d\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, string(c.input), string(exp), string(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
