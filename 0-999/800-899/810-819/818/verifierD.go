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
	n     int
	alice int
	cars  []int
}

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "818D.go")
	bin := filepath.Join(os.TempDir(), "ref818D.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return bin, nil
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &bytes.Buffer{}
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rand.Seed(4)
	tests := make([]Test, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(50) + 1
		cars := make([]int, n)
		for j := range cars {
			cars[j] = rand.Intn(10) + 1
		}
		alice := rand.Intn(10) + 1
		tests[i] = Test{n, alice, cars}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.alice))
		for j, c := range t.cars {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(c))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if want != got {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:%sexpected:%s\nactual:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
