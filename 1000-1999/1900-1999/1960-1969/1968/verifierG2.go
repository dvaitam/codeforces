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

// Embedded correct solver for 1968G2 (z-function + DSU approach)
func solve(n, l, r int, s string) string {
	z := zFunction(s)
	buckets := make([][]int, n+1)
	parent := make([]int, n+2)
	for i := 1; i <= n+1; i++ {
		parent[i] = i
	}

	for i := 1; i < n; i++ {
		pos := i + 1
		if z[i] == 0 {
			parent[pos] = findDSU(parent, pos+1)
		} else {
			buckets[z[i]] = append(buckets[z[i]], pos)
		}
	}

	best := make([]int, n+2)

	for x := 1; x <= n; x++ {
		limit := n - x + 1
		cnt := 0
		pos := 1
		for pos <= limit {
			cnt++
			pos = findDSU(parent, pos+x)
		}
		if x > best[cnt] {
			best[cnt] = x
		}
		for _, p := range buckets[x] {
			parent[p] = findDSU(parent, p+1)
		}
	}

	for k := n - 1; k >= 1; k-- {
		if best[k+1] > best[k] {
			best[k] = best[k+1]
		}
	}

	parts := make([]string, 0, r-l+1)
	for k := l; k <= r; k++ {
		parts = append(parts, fmt.Sprintf("%d", best[k]))
	}
	return strings.Join(parts, " ")
}

func findDSU(parent []int, x int) int {
	root := x
	for parent[root] != root {
		root = parent[root]
	}
	for parent[x] != x {
		p := parent[x]
		parent[x] = root
		x = p
	}
	return root
}

func zFunction(s string) []int {
	n := len(s)
	z := make([]int, n)
	if n == 0 {
		return z
	}
	z[0] = n
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i <= r {
			v := z[i-l]
			if v > r-i+1 {
				v = r - i + 1
			}
			z[i] = v
		}
		for i+z[i] < n && s[z[i]] == s[i+z[i]] {
			z[i]++
		}
		if i+z[i]-1 > r {
			l = i
			r = i + z[i] - 1
		}
	}
	return z
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	l := rng.Intn(n) + 1
	r := l + rng.Intn(n-l+1)
	letters := []byte("abcdefghijklmnopqrstuvwxyz")
	sb := make([]byte, n)
	for i := range sb {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	s := string(sb)
	input := fmt.Sprintf("1\n%d %d %d\n%s\n", n, l, r, s)
	expect := solve(n, l, r, s)
	return input, expect
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		inp, exp := genCase(rng)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d: expected %s got %s\ninput:\n%s", i+1, exp, got, inp)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
