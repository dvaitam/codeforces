package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseF struct {
	n   int
	arr []int64
	exp string
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveF(n int, arr []int64) string {
	full := 1<<n - 1
	memo := make(map[int]int)
	var dfs func(int) int
	dfs = func(mask int) int {
		if mask == full {
			return 0
		}
		if v, ok := memo[mask]; ok {
			return v
		}
		i := 0
		for ; i < n; i++ {
			if mask>>i&1 == 0 {
				break
			}
		}
		best := dfs(mask | 1<<i)
		for j := i + 1; j < n; j++ {
			if mask>>j&1 == 0 && abs64(arr[i]-arr[j]) == 1 {
				val := 1 + dfs(mask|1<<i|1<<j)
				if val > best {
					best = val
				}
			}
		}
		memo[mask] = best
		return best
	}
	pairs := dfs(0)
	ans := n - pairs
	return fmt.Sprint(ans)
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

func generateTests() []testCaseF {
	rng := rand.New(rand.NewSource(7))
	cases := make([]testCaseF, 100)
	for i := range cases {
		n := rng.Intn(6) + 1
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(10))
		}
		cases[i] = testCaseF{n: n, arr: arr, exp: solveF(n, arr)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for j := 0; j < tc.n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(tc.arr[j], 10))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
