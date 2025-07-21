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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveC(ti, td, tr, te int, a, b string) string {
	n, m := len(a), len(b)
	maxCost := ti
	if td > maxCost {
		maxCost = td
	}
	if tr > maxCost {
		maxCost = tr
	}
	if te > maxCost {
		maxCost = te
	}
	maxdist := (n + m) * maxCost
	H := make([][]int, n+2)
	for i := range H {
		H[i] = make([]int, m+2)
	}
	H[0][0] = maxdist
	for i := 0; i <= n; i++ {
		H[i+1][0] = maxdist
		H[i+1][1] = i * td
	}
	for j := 0; j <= m; j++ {
		H[0][j+1] = maxdist
		H[1][j+1] = j * ti
	}
	da := make([]int, 26)
	for i := 1; i <= n; i++ {
		db := 0
		for j := 1; j <= m; j++ {
			i1 := da[b[j-1]-'a']
			j1 := db
			cost := tr
			if a[i-1] == b[j-1] {
				cost = 0
				db = j
			}
			h0 := H[i][j] + cost
			h1 := H[i+1][j] + ti
			h2 := H[i][j+1] + td
			h := min(h0, min(h1, h2))
			trans := H[i1][j1] + (i-i1-1)*td + te + (j-j1-1)*ti
			if trans < h {
				h = trans
			}
			H[i+1][j+1] = h
		}
		da[a[i-1]-'a'] = i
	}
	res := H[n+1][m+1]
	return fmt.Sprintf("%d", res)
}

func randString(rng *rand.Rand, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCaseC(rng *rand.Rand) (string, string) {
	ti := rng.Intn(5) + 1
	td := rng.Intn(5) + 1
	tr := rng.Intn(5) + 1
	te := rng.Intn(5) + 1
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	a := randString(rng, n)
	b := randString(rng, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", ti, td, tr, te)
	fmt.Fprintf(&sb, "%s\n%s\n", a, b)
	expected := solveC(ti, td, tr, te, a, b)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
