package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveB(in string) string {
	reader := bufio.NewReader(strings.NewReader(in))
	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return ""
	}
	var x, y string
	fmt.Fscan(reader, &x)
	fmt.Fscan(reader, &y)
	px, py := len(x), len(y)
	g := gcd(px, py)
	ty := int64(py / g)
	fullCycles := n / ty
	var totalPairsPerCycle int64
	for r := 0; r < g; r++ {
		var cntX [26]int64
		for i := r; i < px; i += g {
			cntX[x[i]-'a']++
		}
		var cntY [26]int64
		for j := r; j < py; j += g {
			cntY[y[j]-'a']++
		}
		for c := 0; c < 26; c++ {
			totalPairsPerCycle += cntX[c] * cntY[c]
		}
	}
	totalMatches := totalPairsPerCycle * fullCycles
	totalPositions := n * int64(px)
	result := totalPositions - totalMatches
	return fmt.Sprintln(result)
}

func genTest(r *rand.Rand) string {
	px := r.Intn(3) + 1
	py := r.Intn(3) + 1
	g := gcd(px, py)
	l := px * py / g
	k := r.Intn(3) + 1
	n := int64(k * l / px)
	m := int64(k * l / py)
	var xb, yb strings.Builder
	for i := 0; i < px; i++ {
		xb.WriteByte(byte('a' + r.Intn(26)))
	}
	for i := 0; i < py; i++ {
		yb.WriteByte(byte('a' + r.Intn(26)))
	}
	return fmt.Sprintf("%d %d\n%s\n%s\n", n, m, xb.String(), yb.String())
}

func runBinary(path, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: verifierB <path-to-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(2))
	const tests = 100
	for i := 0; i < tests; i++ {
		in := genTest(r)
		expect := strings.TrimSpace(solveB(in))
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
