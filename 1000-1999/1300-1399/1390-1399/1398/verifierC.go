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
	t := 1
	_ = sc.Text()
	outputs := make([]string, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		s := sc.Text()
		freq := map[int]int{0: 1}
		sum := 0
		var ans int64
		for i := 0; i < n; i++ {
			sum += int(s[i]-'0') - 1
			if c, ok := freq[sum]; ok {
				ans += int64(c)
			}
			freq[sum]++
		}
		outputs[caseIdx] = fmt.Sprint(ans)
	}
	return strings.Join(outputs, "\n")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	sb := make([]byte, n)
	for i := 0; i < n; i++ {
		sb[i] = byte('0' + rng.Intn(3))
	}
	var b strings.Builder
	b.WriteString("1\n")
	b.WriteString(fmt.Sprintln(n))
	b.WriteString(string(sb))
	b.WriteByte('\n')
	return b.String()
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
