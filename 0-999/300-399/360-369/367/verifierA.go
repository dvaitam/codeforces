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

type query struct{ l, r int }

func solveA(s string, qs []query) []string {
	n := len(s)
	px := make([]int, n+1)
	py := make([]int, n+1)
	pz := make([]int, n+1)
	for i, ch := range s {
		px[i+1] = px[i]
		py[i+1] = py[i]
		pz[i+1] = pz[i]
		switch ch {
		case 'x':
			px[i+1]++
		case 'y':
			py[i+1]++
		case 'z':
			pz[i+1]++
		}
	}
	res := make([]string, len(qs))
	for i, q := range qs {
		l, r := q.l, q.r
		length := r - l + 1
		if length < 3 {
			res[i] = "YES"
			continue
		}
		cx := px[r] - px[l-1]
		cy := py[r] - py[l-1]
		cz := pz[r] - pz[l-1]
		mx := cx
		if cy > mx {
			mx = cy
		}
		if cz > mx {
			mx = cz
		}
		mn := cx
		if cy < mn {
			mn = cy
		}
		if cz < mn {
			mn = cz
		}
		if mx-mn > 1 {
			res[i] = "NO"
		} else {
			res[i] = "YES"
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(50) + 1
	sb := strings.Builder{}
	letters := []byte{'x', 'y', 'z'}
	for i := 0; i < n; i++ {
		sb.WriteByte(letters[rng.Intn(3)])
	}
	s := sb.String()
	m := rng.Intn(20) + 1
	qs := make([]query, m)
	input := fmt.Sprintf("%s\n%d\n", s, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		qs[i] = query{l, r}
		input += fmt.Sprintf("%d %d\n", l, r)
	}
	ans := solveA(s, qs)
	expected := strings.Join(ans, "\n")
	return input, expected
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
