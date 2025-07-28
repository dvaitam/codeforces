package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n   int
	a   []int64
	exp string
}

func canZero(a []int64) bool {
	n := len(a)
	for i := 0; i <= n-3; i++ {
		x := a[i]
		if a[i+1] < 2*x || a[i+2] < x {
			return false
		}
		a[i+1] -= 2 * x
		a[i+2] -= x
	}
	return a[n-2] == 0 && a[n-1] == 0
}

func solveB(arr []int64) string {
	b := make([]int64, len(arr))
	copy(b, arr)
	if canZero(b) {
		return "YES"
	}
	return "NO"
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCaseB {
	rng := rand.New(rand.NewSource(2))
	cases := make([]testCaseB, 100)
	for i := range cases {
		n := rng.Intn(6) + 3
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(20))
		}
		cases[i] = testCaseB{n: n, a: arr, exp: solveB(arr)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintln(&sb, 1)
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.a[j]))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(strings.ToUpper(got)) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
