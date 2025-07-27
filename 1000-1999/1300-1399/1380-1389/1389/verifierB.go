package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestCase struct {
	n, k, z int
	a       []int
}

func solveCase(tc TestCase) int {
	n := tc.n
	k := tc.k
	z := tc.z
	a := tc.a
	best := 0
	var dfs func(pos, moves, leftUsed int, lastLeft bool, score int)
	dfs = func(pos, moves, leftUsed int, lastLeft bool, score int) {
		if moves == k {
			if score > best {
				best = score
			}
			return
		}
		if pos < n-1 {
			dfs(pos+1, moves+1, leftUsed, false, score+a[pos+1])
		}
		if pos > 0 && !lastLeft && leftUsed < z {
			dfs(pos-1, moves+1, leftUsed+1, true, score+a[pos-1])
		}
	}
	dfs(0, 0, 0, false, a[0])
	return best
}

func run(bin, input string) (string, error) {
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

func genTest(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	k := rng.Intn(n-1) + 1
	z := rng.Intn(3)
	if z > k {
		z = k
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10) + 1
	}
	tc := TestCase{n, k, z, a}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, z))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", solveCase(tc))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, exp := genTest(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
