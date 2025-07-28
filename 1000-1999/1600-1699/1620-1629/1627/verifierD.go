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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	arr := make([]int, n)
	maxv := 0
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
		if arr[i] > maxv {
			maxv = arr[i]
		}
	}
	present := make([]bool, maxv+1)
	for _, v := range arr {
		present[v] = true
	}
	cnt := make([]int, maxv+1)
	gval := make([]int, maxv+1)
	for g := 1; g <= maxv; g++ {
		val := 0
		c := 0
		for m := g; m <= maxv; m += g {
			if present[m] {
				c++
				if val == 0 {
					val = m
				} else {
					val = gcd(val, m)
				}
			}
		}
		cnt[g] = c
		gval[g] = val
	}
	ans := 0
	for g := 1; g <= maxv; g++ {
		if cnt[g] >= 2 && gval[g] == g && !present[g] {
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
}

func runProg(bin, input string) (string, error) {
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

func generateTests() []string {
	rng := rand.New(rand.NewSource(4))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 2
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		used := make(map[int]struct{})
		for j := 0; j < n; j++ {
			val := rng.Intn(200) + 1
			for {
				if _, ok := used[val]; !ok {
					break
				}
				val = rng.Intn(200) + 1
			}
			used[val] = struct{}{}
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		expect := solveD(t)
		got, err := runProg(bin, t)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
