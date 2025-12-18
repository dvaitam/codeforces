package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Test struct {
	input string
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "904C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func randWordWith(n int, c byte) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rand.Intn(26))
	}
	pos := rand.Intn(n)
	b[pos] = c
	return string(b)
}

func randWordWithout(n int, c byte) string {
	b := make([]byte, n)
	for i := range b {
		for {
			x := byte('a' + rand.Intn(26))
			if x != c {
				b[i] = x
				break
			}
		}
	}
	return string(b)
}

func genTests() []Test {
	rand.Seed(time.Now().UnixNano())
	tests := make([]Test, 0, 100)
	for i := 0; i < 100; i++ {
		// Choose a secret letter
		secret := byte('a' + rand.Intn(26))
		
		n := rand.Intn(15) + 1 // Total actions including the last one
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		
		// Generate n-1 actions
		for j := 0; j < n-1; j++ {
			t := rand.Intn(3)
			switch t {
			case 0: // . word (must NOT contain secret)
				sb.WriteString(". " + randWordWithout(rand.Intn(6)+1, secret) + "\n")
			case 1: // ! word (MUST contain secret)
				sb.WriteString("! " + randWordWith(rand.Intn(6)+1, secret) + "\n")
			case 2: // ? guess (MUST NOT be secret, as it's not the last action)
				var guess byte
				for {
					guess = byte('a' + rand.Intn(26))
					if guess != secret {
						break
					}
				}
				sb.WriteString("? " + string(guess) + "\n")
			}
		}
		// Last action: ? secret
		sb.WriteString("? " + string(secret) + "\n")
		
		tests = append(tests, Test{sb.String()})
	}
	return tests
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}

	if len(args) != 1 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := args[0]
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
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}