package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func expectedF(input string) string {
	// for k=0 output is always "1\n0"
	return "1\n0"
}

func genTestsF() []string {
	rand.Seed(6)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(3) + 3 // 3..5
		m := rand.Intn(3) + 3
		k := 0
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if rand.Intn(2) == 0 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('1')
				}
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierF.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsF()
	want := strings.TrimSpace(expectedF(""))
	for i, tcase := range tests {
		got, err := runBinary(bin, tcase)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed.\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, tcase, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
