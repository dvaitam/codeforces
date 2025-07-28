package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Case struct{ a []int }

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(8) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(10) + 1
		}
		cases[i] = Case{a: arr}
	}
	return cases
}

func solve(a []int) int {
	n := len(a)
	prefix := make([]int, n+1)
	pos := make(map[int]int, n+1)
	pos[0] = 0
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
		pos[prefix[i+1]] = i + 1
	}
	ans := n
	for i := 1; i <= n; i++ {
		target := prefix[i]
		last := i
		maxLen := i
		valid := true
		for last < n {
			needed := prefix[last] + target
			idx, ok := pos[needed]
			if !ok {
				valid = false
				break
			}
			segLen := idx - last
			if segLen > maxLen {
				maxLen = segLen
			}
			last = idx
		}
		if valid && last == n && maxLen < ans {
			ans = maxLen
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

func runCase(bin string, c Case) error {
	exp := solve(c.a)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(c.a)))
	for i, v := range c.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil || got != exp {
		return fmt.Errorf("expected %d got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
