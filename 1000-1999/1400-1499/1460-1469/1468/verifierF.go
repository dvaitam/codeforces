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
)

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1468F.go")
	bin := filepath.Join(os.TempDir(), "ref1468F.bin")
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

func genCases() []Case {
	rng := rand.New(rand.NewSource(1468))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 1
		var buf strings.Builder
		fmt.Fprintf(&buf, "1\n%d\n", n)
		for j := 0; j < n; j++ {
			x := rng.Intn(10)
			y := rng.Intn(10)
			u := rng.Intn(10)
			v := rng.Intn(10)
			fmt.Fprintf(&buf, "%d %d %d %d\n", x, y, u, v)
		}
		cases[i] = Case{[]byte(buf.String())}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
