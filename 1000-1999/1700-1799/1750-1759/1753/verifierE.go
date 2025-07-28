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

type test struct {
	input    string
	expected string
}

func solveCase(n int, b, pc, mc int64, opsType []byte, opsVal []int64) string {
	orig := int64(1)
	prodAll := int64(1)
	plusSum := int64(0)
	for i := 0; i < n; i++ {
		if opsType[i] == '+' {
			orig += opsVal[i]
			plusSum += opsVal[i]
		} else {
			orig *= opsVal[i]
			if opsVal[i] > 1 {
				prodAll *= opsVal[i]
			}
		}
	}
	sufProd := make([]int64, n+1)
	sufProd[n] = 1
	for i := n - 1; i >= 0; i-- {
		sufProd[i] = sufProd[i+1]
		if opsType[i] == '*' && opsVal[i] > 1 {
			sufProd[i] *= opsVal[i]
		}
	}
	prefMult := int64(0)
	minCost := int64(1 << 60)
	plusSuffix := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		plusSuffix[i] = plusSuffix[i+1]
		if opsType[i] == '+' {
			plusSuffix[i]++
		}
	}
	for i := 0; i <= n; i++ {
		multBefore := prefMult
		plusAfter := int64(plusSuffix[i])
		cost := multBefore*mc + plusAfter*pc
		if cost < minCost {
			minCost = cost
		}
		if i < n && opsType[i] == '*' && opsVal[i] > 1 {
			prefMult++
		}
	}
	if minCost <= b {
		result := (1 + plusSum) * prodAll
		return fmt.Sprintf("%d", result)
	}
	profits := make([]int64, 0)
	for i := 0; i < n; i++ {
		if opsType[i] == '+' {
			gain := opsVal[i] * (prodAll - sufProd[i])
			if gain > 0 {
				profits = append(profits, gain)
			}
		}
	}
	sort.Slice(profits, func(i, j int) bool { return profits[i] > profits[j] })
	moves := int64(len(profits))
	if moves*pc > b {
		moves = b / pc
	}
	sumGain := int64(0)
	for i := int64(0); i < moves; i++ {
		sumGain += profits[i]
	}
	result := orig + sumGain
	return fmt.Sprintf("%d", result)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(46))
	var tests []test
	tests = append(tests, test{input: "1 0 1 1\n", expected: solveCase(1, 0, 1, 1, []byte{'+'}, []int64{0})})
	for len(tests) < 100 {
		n := rng.Intn(4) + 1
		b := int64(rng.Intn(5))
		pc := int64(rng.Intn(3) + 1)
		mc := int64(rng.Intn(3) + 1)
		opsT := make([]byte, n)
		opsV := make([]int64, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, b, pc, mc))
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				opsT[i] = '+'
			} else {
				opsT[i] = '*'
			}
			opsV[i] = int64(rng.Intn(4) + 1)
			sb.WriteByte(opsT[i])
			sb.WriteByte(' ')
			sb.WriteString(fmt.Sprintf("%d\n", opsV[i]))
		}
		tests = append(tests, test{input: sb.String(), expected: solveCase(n, b, pc, mc, opsT, opsV)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
