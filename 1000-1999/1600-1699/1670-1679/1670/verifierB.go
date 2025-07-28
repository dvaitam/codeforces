package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func run(cmdPath string, input string) (string, error) {
	cmd := exec.Command(cmdPath)
	if strings.HasSuffix(cmdPath, ".go") {
		cmd = exec.Command("go", "run", cmdPath)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1670B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []string {
	rand.Seed(2)
	tests := make([]string, 0, 100)
	letters := "abcdefghijklmnopqrstuvwxyz"
	for i := 0; i < 100; i++ {
		n := rand.Intn(20) + 2
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(n) + "\n")
		for j := 0; j < n; j++ {
			sb.WriteByte(letters[rand.Intn(26)])
		}
		sb.WriteByte('\n')
		k := rand.Intn(26) + 1
		sb.WriteString(strconv.Itoa(k))
		sb.WriteByte(' ')
		chosen := rand.Perm(26)[:k]
		for j, idx := range chosen {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteByte(letters[idx])
		}
		sb.WriteByte('\n')
		tests = append(tests, sb.String())
	}
	// edge cases
	sb := strings.Builder{}
	sb.WriteString("1\n2\naa\n1 a\n")
	tests = append(tests, sb.String())
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
		exp, err := run(ref, tc)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, tc, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
