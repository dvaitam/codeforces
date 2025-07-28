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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("python3", bin)
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

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var T int
	fmt.Fscan(in, &T)
	var sb strings.Builder
	for ; T > 0; T-- {
		var n, k int
		var x int64
		fmt.Fscan(in, &n, &k, &x)
		var s string
		fmt.Fscan(in, &s)
		letters := make([]byte, 0, len(s))
		bases := make([]int64, 0)
		for i := 0; i < len(s); {
			if s[i] == 'a' {
				letters = append(letters, 'a')
				i++
			} else {
				j := i
				for j < len(s) && s[j] == '*' {
					j++
				}
				base := int64(j-i)*int64(k) + 1
				bases = append(bases, base)
				letters = append(letters, '*')
				i = j
			}
		}
		x--
		cntB := make([]int64, len(bases))
		for i := len(bases) - 1; i >= 0; i-- {
			cntB[i] = x % bases[i]
			x /= bases[i]
		}
		var out strings.Builder
		idx := 0
		for _, ch := range letters {
			if ch == 'a' {
				out.WriteByte('a')
			} else {
				if cntB[idx] > 0 {
					out.WriteString(strings.Repeat("b", int(cntB[idx])))
				}
				idx++
			}
		}
		sb.WriteString(out.String())
		if T > 1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	k := rng.Intn(3) + 1
	var sBuilder strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			sBuilder.WriteByte('*')
		} else {
			sBuilder.WriteByte('a')
		}
	}
	s := sBuilder.String()
	// compute total combinations
	total := int64(1)
	for i := 0; i < len(s); {
		if s[i] == 'a' {
			i++
			continue
		}
		j := i
		for j < len(s) && s[j] == '*' {
			j++
		}
		base := int64(j-i)*int64(k) + 1
		total *= base
		i = j
	}
	if total <= 0 {
		total = 1
	}
	x := rng.Int63n(total) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d %d\n%s\n", n, k, x, s)
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := generateCase(rng)
		expect := solve(input)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
