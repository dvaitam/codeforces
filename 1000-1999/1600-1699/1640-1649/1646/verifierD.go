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
		return "", fmt.Errorf("read reference source: %v", err)
	}
	bin := os.TempDir() + "/1646D_ref.bin"
	if strings.Contains(string(content), "#include") {
		cppSrc := os.TempDir() + "/1646D_ref.cpp"
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			return "", fmt.Errorf("write cpp source: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", bin, cppSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build c++ reference failed: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", bin, refSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
		}
	}
	return bin, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(45))
	tests := []string{"2\n1 2\n"}
	for len(tests) < 100 {
		n := rng.Intn(15) + 2
		edges := make([][2]int, n-1)
		for i := 2; i <= n; i++ {
			p := rng.Intn(i-1) + 1
			edges[i-2] = [2]int{p, i}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		return
	}
	bin := os.Args[1]
	oracle, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	tests := generateTests()
	for i, input := range tests {
		expect, err := runBinary(oracle, input)
		if err != nil {
			fmt.Printf("oracle failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		// Only compare first line (count and sum); weights may differ for equivalent solutions
		expectFirst := strings.SplitN(strings.TrimSpace(expect), "\n", 2)[0]
		gotFirst := strings.SplitN(strings.TrimSpace(got), "\n", 2)[0]
		if strings.TrimSpace(gotFirst) != strings.TrimSpace(expectFirst) {
			fmt.Printf("wrong answer on test %d\ninput:%sexpected:%s\ngot:%s\n", i+1, input, expectFirst, gotFirst)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
