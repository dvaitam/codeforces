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

const MOD int = 1_000_000_007

func weight(a []int, idx []int) int {
	if len(idx) == 0 {
		return 0
	}
	last := idx[len(idx)-1]
	maxAfter := 0
	for i := last + 1; i < len(a); i++ {
		if a[i] > maxAfter {
			maxAfter = a[i]
		}
	}
	w := 0
	for _, p := range idx {
		if a[p] < maxAfter {
			w++
		}
	}
	return w % MOD
}

func dfs(a []int, start int, subseq []int, ans *int) {
	if len(subseq) > 0 {
		*ans = (*ans + weight(a, subseq)) % MOD
	}
	for i := start; i < len(a); i++ {
		if len(subseq) == 0 || a[i] > a[subseq[len(subseq)-1]] {
			subseq = append(subseq, i)
			dfs(a, i+1, subseq, ans)
			subseq = subseq[:len(subseq)-1]
		}
	}
}

func solveCaseG(a []int) int {
	ans := 0
	dfs(a, 0, []int{}, &ans)
	return ans % MOD
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseG(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(10) + 1
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i, v := range a {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	exp := fmt.Sprintf("%d", solveCaseG(a))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseG(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
