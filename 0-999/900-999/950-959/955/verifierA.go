package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	hh, mm     int
	H, D, C, N int
}

func (t Test) Input() string {
	return fmt.Sprintf("%02d %02d\n%d %d %d %d\n", t.hh, t.mm, t.H, t.D, t.C, t.N)
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "955A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(0))
	tests := make([]Test, 0, 110)
	for i := 0; i < 100; i++ {
		tests = append(tests, Test{
			hh: rng.Intn(24),
			mm: rng.Intn(60),
			H:  rng.Intn(500) + 1,
			D:  rng.Intn(10) + 1,
			C:  rng.Intn(200) + 1,
			N:  rng.Intn(10) + 1,
		})
	}
	// edge cases around discount time
	tests = append(tests,
		Test{19, 59, 10, 1, 20, 5},
		Test{20, 0, 10, 2, 15, 3},
		Test{23, 30, 1, 1, 1, 1},
		Test{0, 0, 100, 5, 100, 10},
		Test{19, 0, 200, 3, 50, 5},
	)
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	candidate := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n" + input)
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
