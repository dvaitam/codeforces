package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func buildLCMs(limit uint64) []uint64 {
	lcms := []uint64{1}
	for i := 1; ; i++ {
		prev := lcms[len(lcms)-1]
		g := gcd(prev, uint64(i))
		l := prev / g * uint64(i)
		if l > limit {
			break
		}
		lcms = append(lcms, l)
	}
	return lcms
}

func expectedAnswer(n uint64, lcms []uint64) int {
	ans := 0
	for i := 1; i < len(lcms); i++ {
		if n%lcms[i] == 0 {
			ans = i
		} else {
			break
		}
	}
	return ans
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	lcms := buildLCMs(1e18)

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n64, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid number on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		exp := expectedAnswer(n64, lcms)
		input := fmt.Sprintf("1\n%d\n", n64)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		want := fmt.Sprintf("%d", exp)
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
