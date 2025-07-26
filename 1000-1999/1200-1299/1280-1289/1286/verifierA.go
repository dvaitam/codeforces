package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n   int
	arr []int
	exp string
}

func solveA(n int, arr []int) string {
	oddExisting, evenExisting := 0, 0
	for _, v := range arr {
		if v != 0 {
			if v%2 == 1 {
				oddExisting++
			} else {
				evenExisting++
			}
		}
	}
	totOdd := (n + 1) / 2
	totEven := n / 2
	oddMissing := totOdd - oddExisting
	evenMissing := totEven - evenExisting

	prefixZero := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefixZero[i] = prefixZero[i-1]
		if arr[i-1] == 0 {
			prefixZero[i]++
		}
	}
	const INF = int(1 << 30)
	dp := make([][][]int, n+1)
	for i := range dp {
		dp[i] = make([][]int, oddMissing+1)
		for j := range dp[i] {
			dp[i][j] = []int{INF, INF}
		}
	}
	dp[0][0][0] = 0
	dp[0][0][1] = 0
	for i := 0; i < n; i++ {
		for o := 0; o <= oddMissing; o++ {
			for last := 0; last < 2; last++ {
				cur := dp[i][o][last]
				if cur == INF {
					continue
				}
				zeroSoFar := prefixZero[i]
				evenUsed := zeroSoFar - o
				if arr[i] != 0 {
					p := arr[i] % 2
					add := 0
					if i > 0 && last != p {
						add = 1
					}
					if dp[i+1][o][p] > cur+add {
						dp[i+1][o][p] = cur + add
					}
				} else {
					if o < oddMissing {
						add := 0
						if i > 0 && last != 1 {
							add = 1
						}
						if dp[i+1][o+1][1] > cur+add {
							dp[i+1][o+1][1] = cur + add
						}
					}
					if evenUsed < evenMissing {
						add := 0
						if i > 0 && last != 0 {
							add = 1
						}
						if dp[i+1][o][0] > cur+add {
							dp[i+1][o][0] = cur + add
						}
					}
				}
			}
		}
	}
	res := dp[n][oddMissing][0]
	if dp[n][oddMissing][1] < res {
		res = dp[n][oddMissing][1]
	}
	return fmt.Sprint(res)
}

func run(bin string, input string) (string, error) {
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

func generateTests() []testCaseA {
	rng := rand.New(rand.NewSource(1))
	cases := make([]testCaseA, 100)
	for i := range cases {
		n := rng.Intn(20) + 1
		perm := rng.Perm(n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = perm[j] + 1
		}
		for j := 0; j < n; j++ {
			if rng.Float64() < 0.3 {
				arr[j] = 0
			}
		}
		cases[i] = testCaseA{n: n, arr: arr, exp: solveA(n, arr)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
			sb.WriteString(fmt.Sprint(tc.arr[j]))
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
