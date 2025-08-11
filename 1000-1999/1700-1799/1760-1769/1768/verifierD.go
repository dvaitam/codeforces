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
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func solveD(n int, p []int) int {
	visited := make([]bool, n+1)
	cycleID := make([]int, n+1)
	cid := 0
	for i := 1; i <= n; i++ {
		if !visited[i] {
			cid++
			j := i
			for !visited[j] {
				visited[j] = true
				cycleID[j] = cid
				j = p[j]
			}
		}
	}
	base := n - cid
	same := false
	for i := 1; i < n; i++ {
		if cycleID[i] == cycleID[i+1] {
			same = true
			break
		}
	}
	ans := base + 1
	if same {
		ans = base - 1
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n int
		p []int
	}

	var cases []test
	cases = append(cases, test{n: 2, p: []int{0, 1, 2}})
	cases = append(cases, test{n: 2, p: []int{0, 2, 1}})
	cases = append(cases, test{n: 3, p: []int{0, 2, 3, 1}})

	for len(cases) < 100 {
		n := rng.Intn(10) + 2
		perm := rand.Perm(n)
		p := make([]int, n+1)
		for i, v := range perm {
			p[i+1] = v + 1
		}
		cases = append(cases, test{n: n, p: p})
	}

	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j := 1; j <= tc.n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.p[j]))
		}
		sb.WriteByte('\n')
		expected := solveD(tc.n, tc.p)
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != 1 {
			fmt.Fprintf(os.Stderr, "case %d: expected single integer got %q\n", i+1, out)
			os.Exit(1)
		}
		val, err := strconv.Atoi(fields[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse integer\n", i+1)
			os.Exit(1)
		}
		if val != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expected, val, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
