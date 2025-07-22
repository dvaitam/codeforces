package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveB(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var pStr, qStr string
	if _, err := fmt.Fscan(reader, &pStr, &qStr); err != nil {
		return ""
	}
	p := new(big.Int)
	p.SetString(pStr, 10)
	q := new(big.Int)
	q.SetString(qStr, 10)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		var aiStr string
		fmt.Fscan(reader, &aiStr)
		ai := new(big.Int)
		ai.SetString(aiStr, 10)
		a[i] = ai
	}
	num := new(big.Int).Set(a[n-1])
	den := big.NewInt(1)
	for i := n - 2; i >= 0; i-- {
		newNum := new(big.Int).Mul(a[i], num)
		newNum.Add(newNum, den)
		newDen := new(big.Int).Set(num)
		num, den = newNum, newDen
	}
	lhs := new(big.Int).Mul(p, den)
	rhs := new(big.Int).Mul(q, num)
	if lhs.Cmp(rhs) == 0 {
		return "YES"
	}
	return "NO"
}

func fractionFromA(a []*big.Int) (*big.Int, *big.Int) {
	n := len(a)
	num := new(big.Int).Set(a[n-1])
	den := big.NewInt(1)
	for i := n - 2; i >= 0; i-- {
		newNum := new(big.Int).Mul(a[i], num)
		newNum.Add(newNum, den)
		newDen := new(big.Int).Set(num)
		num, den = newNum, newDen
	}
	return num, den
}

func genTestB() string {
	n := rand.Intn(10) + 1
	a := make([]*big.Int, n)
	for i := 0; i < n; i++ {
		a[i] = big.NewInt(int64(rand.Intn(1000) + 1))
	}
	num, den := fractionFromA(a)
	// decide if we want equal or not
	equal := rand.Intn(2) == 0
	p := new(big.Int).Set(num)
	q := new(big.Int).Set(den)
	if !equal {
		// modify p slightly
		p.Add(p, big.NewInt(int64(rand.Intn(5)+1)))
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %s\n", p.String(), q.String())
	fmt.Fprintln(&buf, n)
	for i, v := range a {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(v.String())
	}
	buf.WriteByte('\n')
	return buf.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		in := genTestB()
		expected := solveB(in)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\noutput: %s\n", i, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
