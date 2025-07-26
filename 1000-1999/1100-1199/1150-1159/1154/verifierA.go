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

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1154A.go")
	bin := filepath.Join(dir, "refA.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return bin, nil
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(1))
	tests := make([]Test, 0, 102)
	for i := 0; i < 100; i++ {
		a := rng.Intn(1000) + 1
		b := rng.Intn(1000) + 1
		c := rng.Intn(1000) + 1
		sums := []int{a + b, a + c, b + c, a + b + c}
		rng.Shuffle(len(sums), func(i, j int) { sums[i], sums[j] = sums[j], sums[i] })
		input := fmt.Sprintf("%d %d %d %d\n", sums[0], sums[1], sums[2], sums[3])
		tests = append(tests, Test{input})
	}
	tests = append(tests, Test{"2 2 2 3\n"})     // a=1,b=1,c=1
	tests = append(tests, Test{"60 30 50 40\n"}) // a=10,b=20,c=30
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
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
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
