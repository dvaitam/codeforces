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
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	content, err := os.ReadFile(refSrc)
	if err != nil {
		return "", fmt.Errorf("failed to read reference source: %v", err)
	}
	tmp := "./refF_1942.bin"
	if strings.Contains(string(content), "#include") {
		cppFile := refSrc + ".cpp"
		if err := os.WriteFile(cppFile, content, 0644); err != nil {
			return "", fmt.Errorf("failed to write cpp file: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", tmp, cppFile)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", tmp, refSrc)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", err
		}
	}
	return tmp, nil
}

func runProg(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stderr = os.Stderr
	cmd.Stdout = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	q := rng.Intn(4) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&buf, "%d", rng.Int63n(1000000000000000000)+1)
		if i+1 < n {
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte('\n')
	for i := 0; i < q; i++ {
		k := rng.Intn(n) + 1
		x := rng.Int63n(1000000000000000000) + 1
		fmt.Fprintf(&buf, "%d %d\n", k, x)
	}
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(48))
	for i := 0; i < 100; i++ {
		test := genTest(rng)
		expected, err := runProg(ref, test)
		if err != nil {
			fmt.Fprintln(os.Stderr, "reference error:", err)
			os.Exit(1)
		}
		got, err := runProg(target, test)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d execution error: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, test, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
