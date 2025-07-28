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

const mod = 998244353

func reachableMasks(n, k int) map[int]struct{} {
	start := (1 << n) - 1
	visited := make(map[int]struct{})
	var dfs func(mask int)
	dfs = func(mask int) {
		if _, ok := visited[mask]; ok {
			return
		}
		visited[mask] = struct{}{}
		// collect positions
		pos := make([]int, 0, n)
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				pos = append(pos, i)
			}
		}
		m := len(pos)
		if m < 2*k+1 {
			return
		}
		// choose center index among existing positions
		for centerIdx := k; centerIdx < m-k; centerIdx++ {
			leftIdx := pos[:centerIdx]
			rightIdx := pos[centerIdx+1:]
			leftComb := make([]int, 0, k)
			var chooseLeft func(start int)
			chooseLeft = func(start int) {
				if len(leftComb) == k {
					rightComb := make([]int, 0, k)
					var chooseRight func(idx int)
					chooseRight = func(idx int) {
						if len(rightComb) == k {
							newMask := mask
							for _, li := range leftComb {
								newMask &^= 1 << li
							}
							for _, ri := range rightComb {
								newMask &^= 1 << ri
							}
							dfs(newMask)
							return
						}
						for j := idx; j < len(rightIdx); j++ {
							rightComb = append(rightComb, rightIdx[j])
							chooseRight(j + 1)
							rightComb = rightComb[:len(rightComb)-1]
						}
					}
					chooseRight(0)
					return
				}
				for i := start; i < len(leftIdx); i++ {
					leftComb = append(leftComb, leftIdx[i])
					chooseLeft(i + 1)
					leftComb = leftComb[:len(leftComb)-1]
				}
			}
			chooseLeft(0)
		}
	}
	dfs(start)
	return visited
}

func bruteE(n int) []int {
	limit := (n - 1) / 2
	res := make([]int, limit)
	for k := 1; k <= limit; k++ {
		masks := reachableMasks(n, k)
		res[k-1] = len(masks) % mod
	}
	return res
}

func generateCaseE(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 3
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	input := sb.String()
	ans := bruteE(n)
	var exp strings.Builder
	for i, v := range ans {
		if i > 0 {
			exp.WriteByte(' ')
		}
		exp.WriteString(fmt.Sprintf("%d", v))
	}
	exp.WriteByte('\n')
	return input, exp.String()
}

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		got = strings.TrimSpace(got)
		if got != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, strings.TrimSpace(exp), got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
