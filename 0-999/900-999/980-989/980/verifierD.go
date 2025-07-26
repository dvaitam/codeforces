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

type Test struct {
	input  string
	output string
}

type Pair struct {
	sign int
	val  int
}

func squareFree(x int) int {
	if x < 0 {
		x = -x
	}
	res := 1
	for p := 2; p*p <= x; p++ {
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt ^= 1
		}
		if cnt == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func solveD(arr []int) string {
	n := len(arr)
	keys := make([]Pair, n)
	for i, v := range arr {
		if v == 0 {
			keys[i] = Pair{0, 0}
			continue
		}
		sign := 1
		if v < 0 {
			sign = -1
		}
		sf := squareFree(v)
		keys[i] = Pair{sign, sf}
	}
	ans := make([]int, n+1)
	for l := 0; l < n; l++ {
		mp := make(map[Pair]int)
		groups := 0
		for r := l; r < n; r++ {
			k := keys[r]
			if k.sign != 0 {
				if mp[k] == 0 {
					groups++
				}
				mp[k]++
			}
			g := groups
			if g == 0 {
				g = 1
			}
			ans[g]++
		}
	}
	out := make([]string, n)
	for i := 1; i <= n; i++ {
		out[i-1] = strconv.Itoa(ans[i])
	}
	return strings.Join(out, " ") + "\n"
}

func generateTests() []Test {
	rand.Seed(42)
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(6) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(11) - 5
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n) + "\n")
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		output := solveD(arr)
		tests = append(tests, Test{input: input, output: output})
	}
	return tests
}

func run(binary string, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := run(binary, t.input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.output) {
			fmt.Printf("Test %d failed. Input: %q\nExpected: %qGot: %q\n", i+1, t.input, t.output, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
