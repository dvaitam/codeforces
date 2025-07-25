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

func solveE(n int, l, r []int) (bool, string) {
	dp := make([][]bool, n+2)
	choice := make([][]int, n+2)
	for i := 0; i <= n+1; i++ {
		dp[i] = make([]bool, n+2)
		choice[i] = make([]int, n+2)
	}
	for i := 1; i <= n+1; i++ {
		dp[i][i-1] = true
	}
	for length := 1; length <= n; length++ {
		for i := 1; i+length-1 <= n; i++ {
			j := i + length - 1
			tmin := l[i] / 2
			tmax := (r[i] - 1) / 2
			if tmin < 0 {
				tmin = 0
			}
			if tmax > j-i {
				tmax = j - i
			}
			for t := tmin; t <= tmax; t++ {
				if dp[i+1][i+t] && dp[i+t+1][j] {
					dp[i][j] = true
					choice[i][j] = t
					break
				}
			}
		}
	}
	if !dp[1][n] {
		return false, ""
	}
	buf := make([]byte, 0, 2*n)
	var build func(i, j int)
	build = func(i, j int) {
		if i > j {
			return
		}
		t := choice[i][j]
		buf = append(buf, '(')
		build(i+1, i+t)
		buf = append(buf, ')')
		build(i+t+1, j)
	}
	build(1, n)
	return true, string(buf)
}

func checkSequence(seq string, n int, l, r []int) bool {
	if len(seq) != 2*n {
		return false
	}
	posStack := []int{}
	idxStack := []int{}
	openID := 0
	for i := 0; i < len(seq); i++ {
		ch := seq[i]
		if ch == '(' {
			openID++
			if openID > n {
				return false
			}
			posStack = append(posStack, i)
			idxStack = append(idxStack, openID)
		} else if ch == ')' {
			if len(posStack) == 0 {
				return false
			}
			pos := posStack[len(posStack)-1]
			idx := idxStack[len(idxStack)-1]
			posStack = posStack[:len(posStack)-1]
			idxStack = idxStack[:len(idxStack)-1]
			dist := i - pos
			if dist < l[idx] || dist > r[idx] {
				return false
			}
		} else {
			return false
		}
	}
	return len(posStack) == 0 && openID == n
}

func generateCase(rng *rand.Rand) (string, []int, []int, bool) {
	n := rng.Intn(6) + 1
	l := make([]int, n+1)
	r := make([]int, n+1)
	for i := 1; i <= n; i++ {
		li := rng.Intn(2*n-1) + 1
		ri := li + rng.Intn(2*n-li)
		l[i] = li
		r[i] = ri
	}
	ok, _ := solveE(n, l, r)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
	}
	return sb.String(), l, r, ok
}

func runCase(bin, input string, l, r []int, expected bool) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	output := strings.TrimSpace(out.String())
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := scanner.Text()
	if !expected {
		if first != "IMPOSSIBLE" {
			return fmt.Errorf("expected IMPOSSIBLE got %s", first)
		}
		return nil
	}
	if first == "IMPOSSIBLE" {
		return fmt.Errorf("expected sequence got IMPOSSIBLE")
	}
	seq := first
	if !checkSequence(seq, len(l)-1, l, r) {
		return fmt.Errorf("invalid sequence %s", seq)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, l, r, ok := generateCase(rng)
		if err := runCase(bin, input, l, r, ok); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
