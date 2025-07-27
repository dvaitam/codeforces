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

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func expected(n, k, c int, arr []int) string {
	maxVal := 1 << c
	counts := make([]int, maxVal)
	var dfs func(int)
	dfs = func(step int) {
		if step == k {
			xor := 0
			for _, v := range arr {
				xor ^= v
			}
			counts[xor]++
			return
		}
		for i := 0; i < n; i++ {
			arr[i]--
			dfs(step + 1)
			arr[i]++
		}
	}
	dfs(0)
	denom := modPow(int64(n), int64(k))
	invDenom := modPow(denom, mod-2)
	res := make([]string, maxVal)
	for x := 0; x < maxVal; x++ {
		val := int64(counts[x]) * invDenom % mod
		res[x] = fmt.Sprint(val)
	}
	return strings.Join(res, " ")
}

func runProgram(bin, input string) (string, error) {
	if _, err := os.Stat(bin); err == nil && !strings.Contains(bin, "/") {
		bin = "./" + bin
	}
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomCase(rng *rand.Rand) (string, int, int, int, []int) {
	c := rng.Intn(3) + 1
	maxVal := 1 << c
	n := rng.Intn(3) + 1
	k := rng.Intn(3) + 1
	for n > maxVal-k {
		n = rng.Intn(3) + 1
	}
	arr := make([]int, n)
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		v := rng.Intn(maxVal-k+1) + k
		for used[v] {
			v = rng.Intn(maxVal-k+1) + k
		}
		used[v] = true
		arr[i] = v
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, c))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String(), n, k, c, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, n, k, c, arr := randomCase(rng)
		want := expected(n, k, c, append([]int{}, arr...))
		out, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t+1, want, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
