package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveD(v []int) int {
	const INF = int(1e18)
	mn := INF
	mx := -INF
	for i := 0; i+1 < len(v); i++ {
		if v[i] < v[i+1] {
			x := (v[i] + v[i+1]) / 2
			if x < mn {
				mn = x
			}
		}
		if v[i] > v[i+1] {
			x := (v[i] + v[i+1] + 1) / 2
			if x > mx {
				mx = x
			}
		}
	}
	if mx == -INF {
		return 0
	} else if mn == INF {
		return v[0]
	} else {
		if mx <= mn {
			return mx
		}
		return -1
	}
}

func expectedD(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	var t int
	fmt.Sscanf(lines[0], "%d", &t)
	idx := 1
	var out strings.Builder
	for i := 0; i < t; i++ {
		var n int
		fmt.Sscanf(lines[idx], "%d", &n)
		idx++
		arrStr := strings.Fields(lines[idx])
		idx++
		v := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Sscanf(arrStr[j], "%d", &v[j])
		}
		out.WriteString(fmt.Sprintf("%d", solveD(v)))
		if i+1 < t {
			out.WriteByte('\n')
		}
	}
	return out.String()
}

func genTestsD() []string {
	rand.Seed(4)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		t := rand.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t))
		for i := 0; i < t; i++ {
			n := rand.Intn(10) + 2
			sb.WriteString(fmt.Sprintf("%d\n", n))
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", rand.Intn(100)+1))
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
		fmt.Fprintf(os.Stderr, "Usage: go run verifierD.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsD()
	for i, tcase := range tests {
		want := expectedD(tcase)
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
