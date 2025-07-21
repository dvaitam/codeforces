package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveC(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	P := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		P[i] = P[i-1] + a[i]
	}
	minP := make([]int64, n+1)
	minP[0] = P[0]
	for i := 1; i <= n; i++ {
		if P[i] < minP[i-1] {
			minP[i] = P[i]
		} else {
			minP[i] = minP[i-1]
		}
	}
	Suf := make([]int64, n+2)
	for i := n; i >= 1; i-- {
		Suf[i] = Suf[i+1] + a[i]
	}
	total := P[n]
	const INF int64 = 1<<63 - 1
	tcost := INF
	for j := 1; j <= n+1; j++ {
		cost := minP[j-1] + Suf[j]
		if cost < tcost {
			tcost = cost
		}
	}
	for i := 1; i <= n; i++ {
		cost := minP[i-1] + Suf[i+1]
		if cost < tcost {
			tcost = cost
		}
	}
	ans := total - 2*tcost
	return fmt.Sprintf("%d\n", ans)
}

func genTestC() (string, string) {
	n := rand.Intn(8) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		v := rand.Intn(21) - 10
		fmt.Fprintf(&sb, "%d ", v)
	}
	sb.WriteString("\n")
	in := sb.String()
	out := solveC(in)
	return in, out
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTestC()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
