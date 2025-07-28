package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func isBeautiful(a, b, c, d int) bool {
	return a < b && c < d && a < c && b < d
}

func expectedB(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var t int
	fmt.Sscanf(lines[0], "%d", &t)
	var out strings.Builder
	idx := 1
	for i := 0; i < t; i++ {
		var a, b, c, d int
		fmt.Sscanf(lines[idx], "%d %d", &a, &b)
		fmt.Sscanf(lines[idx+1], "%d %d", &c, &d)
		idx += 2
		for r := 0; r < 4; r++ {
			if isBeautiful(a, b, c, d) {
				out.WriteString("YES\n")
				break
			}
			a, b, c, d = c, a, d, b
			if r == 3 {
				out.WriteString("NO\n")
			}
		}
	}
	return strings.TrimSpace(out.String())
}

func genTestsB() []string {
	rand.Seed(2)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		t := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			nums := rand.Perm(100)[:4]
			sb.WriteString(fmt.Sprintf("%d %d\n%d %d\n", nums[0]+1, nums[1]+1, nums[2]+1, nums[3]+1))
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
		fmt.Fprintf(os.Stderr, "Usage: go run verifierB.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, t := range tests {
		want := expectedB(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("Test %d failed.\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, t, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
