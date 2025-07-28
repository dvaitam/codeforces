package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const mod = 1_000_000_007

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func isBad(arr []int) bool {
	n := len(arr)
	tmp := append([]int(nil), arr...)
	sort.Ints(tmp)
	return tmp[n/2] == arr[n/2]
}

func okPrefix(p []int, pos int) bool {
	for l := 3; l <= pos+1; l += 2 {
		start := pos - l + 1
		sub := p[start : pos+1]
		for _, v := range sub {
			if v == -1 {
				goto next
			}
		}
		if isBad(sub) {
			return false
		}
	next:
	}
	return true
}

func count(n int, fixed []int) int {
	used := make([]bool, n+1)
	for _, v := range fixed {
		if v != -1 {
			used[v] = true
		}
	}
	perm := make([]int, n)
	copy(perm, fixed)
	ans := 0
	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			ans = (ans + 1) % mod
			return
		}
		if perm[pos] != -1 {
			if okPrefix(perm, pos) {
				dfs(pos + 1)
			}
			return
		}
		for v := 1; v <= n; v++ {
			if used[v] {
				continue
			}
			perm[pos] = v
			used[v] = true
			if okPrefix(perm, pos) {
				dfs(pos + 1)
			}
			used[v] = false
			perm[pos] = -1
		}
	}
	dfs(0)
	return ans
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	arr := make([]int, n)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			arr[i] = -1
		} else {
			for {
				v := rng.Intn(n) + 1
				if !used[v] {
					used[v] = true
					arr[i] = v
					break
				}
			}
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", count(n, arr))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
