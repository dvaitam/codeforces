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
	input string
}

func buildRef() (string, error) {
	ref := "refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1200A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func runExe(path, input string) (string, error) {
	if !strings.Contains(path, "/") {
		path = "./" + path
	}
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTests() []Test {
	rng := rand.New(rand.NewSource(0))
	tests := make([]Test, 0, 105)
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		rooms := make([]bool, 10)
		for j := 0; j < n; j++ {
			ops := make([]byte, 0, 12)
			free := 0
			for _, v := range rooms {
				if !v {
					free++
				}
			}
			if free > 0 {
				ops = append(ops, 'L', 'R')
			}
			for k := 0; k < 10; k++ {
				if rooms[k] {
					ops = append(ops, byte('0'+k))
				}
			}
			ch := ops[rng.Intn(len(ops))]
			sb.WriteByte(ch)
			switch ch {
			case 'L':
				for idx := 0; idx < 10; idx++ {
					if !rooms[idx] {
						rooms[idx] = true
						break
					}
				}
			case 'R':
				for idx := 9; idx >= 0; idx-- {
					if !rooms[idx] {
						rooms[idx] = true
						break
					}
				}
			default:
				rooms[int(ch-'0')] = false
			}
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests, Test{"1\nL\n"})
	tests = append(tests, Test{"2\nLR\n"})
	tests = append(tests, Test{"3\nLR0\n"})
	tests = append(tests, Test{"10\nLLLLLLLLLL\n"})
	tests = append(tests, Test{"10\nRRRRRRRRRR\n"})
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
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
