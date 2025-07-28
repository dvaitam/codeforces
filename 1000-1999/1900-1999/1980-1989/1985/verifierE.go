package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1985E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func buildCase(x, y, z, k int64) []byte {
	return []byte(fmt.Sprintf("1\n%d %d %d %d\n", x, y, z, k))
}

func genRandomCase(rng *rand.Rand) []byte {
	x := int64(rng.Intn(5) + 1)
	y := int64(rng.Intn(5) + 1)
	z := int64(rng.Intn(5) + 1)
	k := int64(rng.Intn(int(x*y*z)) + 1)
	return buildCase(x, y, z, k)
}

func genTests() [][]byte {
	rng := rand.New(rand.NewSource(5))
	tests := [][]byte{
		buildCase(1, 1, 1, 1),
		buildCase(2, 2, 2, 4),
		buildCase(3, 3, 3, 27),
	}
	for len(tests) < 100 {
		tests = append(tests, genRandomCase(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
		exp, err := runExe(ref, tc)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n%s", i+1, err, exp)
			os.Exit(1)
		}
		got, err := runExe(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, string(tc), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
