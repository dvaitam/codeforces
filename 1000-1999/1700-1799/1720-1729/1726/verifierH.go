package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type pt struct{ x, y int }

type Test struct {
	pts []pt
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(t.pts)))
	for _, p := range t.pts {
		sb.WriteString(fmt.Sprintf("%d %d\n", p.x, p.y))
	}
	return sb.String()
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
	ref := "./refH.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1726H.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(7)
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		w := rand.Intn(5) + 1
		h := rand.Intn(5) + 1
		sx := rand.Intn(10) - 5
		sy := rand.Intn(10) - 5
		pts := []pt{{sx, sy}, {sx + w, sy}, {sx + w, sy + h}, {sx, sy + h}}
		tests = append(tests, Test{pts: pts})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierH.go /path/to/binary")
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
		input := tc.Input()
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:%sExpected:%sGot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
