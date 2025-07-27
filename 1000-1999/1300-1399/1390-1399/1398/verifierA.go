package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func solveCase(input string) string {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	sc.Scan() // t
	t, _ := strconv.Atoi(sc.Text())
	outLines := make([]string, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			sc.Scan()
			v, _ := strconv.Atoi(sc.Text())
			arr[i] = v
		}
		ans := "-1"
		for i := 2; i < n; i++ {
			if arr[0]+arr[1] <= arr[i] {
				ans = fmt.Sprintf("%d %d %d", 1, 2, i+1)
				break
			}
		}
		outLines[caseIdx] = ans
	}
	return strings.Join(outLines, "\n")
}

func generateCase(rng *rand.Rand) string {
	t := 1
	n := rng.Intn(8) + 3
	a := make([]int, n)
	a[0] = rng.Intn(5) + 1
	a[1] = a[0] + rng.Intn(5)
	for i := 2; i < n; i++ {
		a[i] = a[i-1] + rng.Intn(5)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(t))
	sb.WriteString(fmt.Sprintln(n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input, expected string) error {
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solveCase(in)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
