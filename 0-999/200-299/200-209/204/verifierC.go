package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	a string
	b string
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveC(a, b string) string {
	n := len(a)
	v1 := make([][]int, 26)
	v2 := make([][]int, 26)
	for i := 0; i < n; i++ {
		c1 := int(a[i] - 'A')
		c2 := int(b[i] - 'A')
		if c1 >= 0 && c1 < 26 {
			v1[c1] = append(v1[c1], i)
		}
		if c2 >= 0 && c2 < 26 {
			v2[c2] = append(v2[c2], i)
		}
	}
	sum1 := make([][]float64, 26)
	sum2 := make([][]float64, 26)
	for i := 0; i < 26; i++ {
		sz := len(v2[i])
		if sz == 0 {
			continue
		}
		sum1[i] = make([]float64, sz)
		sum2[i] = make([]float64, sz)
		for j := 0; j < sz; j++ {
			val := float64(v2[i][j] + 1)
			if j == 0 {
				sum1[i][j] = val
			} else {
				sum1[i][j] = sum1[i][j-1] + val
			}
		}
		for j := sz - 1; j >= 0; j-- {
			val := float64(n - v2[i][j])
			if j == sz-1 {
				sum2[i][j] = val
			} else {
				sum2[i][j] = sum2[i][j+1] + val
			}
		}
	}
	var res float64
	for i := 0; i < 26; i++ {
		if len(v2[i]) == 0 {
			continue
		}
		for _, pos := range v1[i] {
			idx := sort.SearchInts(v2[i], pos)
			if idx > 0 {
				res += float64(n-pos) * sum1[i][idx-1]
			}
			if idx < len(v2[i]) {
				res += float64(pos+1) * sum2[i][idx]
			}
		}
	}
	div := float64(n) * float64(n+1) * float64(2*n+1) / 6.0
	return fmt.Sprintf("%.10f", res/div)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(3))
	tests := make([]testCase, 0, 100)
	fixed := []testCase{{a: "AB", b: "BA"}, {a: "AAAA", b: "BBBB"}}
	tests = append(tests, fixed...)
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		sb1 := make([]byte, n)
		sb2 := make([]byte, n)
		for i := 0; i < n; i++ {
			sb1[i] = letters[rng.Intn(26)]
			sb2[i] = letters[rng.Intn(26)]
		}
		tests = append(tests, testCase{a: string(sb1), b: string(sb2)})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d\n%s\n%s\n", len(t.a), t.a, t.b)
		expect := strings.TrimSpace(solveC(t.a, t.b))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
