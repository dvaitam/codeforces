package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveE(arr []int) string {
	n := len(arr)
	asc, desc := 0, 0
	for i := 0; i < n; i++ {
		if arr[i] == i+1 {
			asc++
		}
		if arr[i] == n-i {
			desc++
		}
	}
	if asc > desc {
		return "First"
	} else if desc > asc {
		return "Second"
	}
	return "Tie"
}

func expectedE(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var t int
	fmt.Sscanf(lines[0], "%d", &t)
	idx := 1
	var out strings.Builder
	for i := 0; i < t; i++ {
		var n int
		fmt.Sscanf(lines[idx], "%d", &n)
		idx++
		arr := make([]int, n)
		fields := strings.Fields(lines[idx])
		for j := 0; j < n; j++ {
			fmt.Sscanf(fields[j], "%d", &arr[j])
		}
		idx++
		out.WriteString(solveE(arr))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func genTestsE() []string {
	rand.Seed(5)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		t := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			n := rand.Intn(10) + 3
			sb.WriteString(fmt.Sprintf("%d\n", n))
			perm := rand.Perm(n)
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", perm[j]+1))
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
		fmt.Fprintf(os.Stderr, "Usage: go run verifierE.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for i, tcase := range tests {
		want := expectedE(tcase)
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
