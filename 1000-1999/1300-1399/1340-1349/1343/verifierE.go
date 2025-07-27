package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

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

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1343E.go")
	return runProg(ref, input)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 2 // 2..7
	maxM := n * (n - 1) / 2
	m := rng.Intn(maxM-(n-1)+1) + (n - 1)
	a := rng.Intn(n) + 1
	b := rng.Intn(n) + 1
	c := rng.Intn(n) + 1
	edges := make([][2]int, 0, m)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{i, p})
	}
	has := func(x, y int) bool {
		for _, e := range edges {
			if (e[0] == x && e[1] == y) || (e[0] == y && e[1] == x) {
				return true
			}
		}
		return false
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v || has(u, v) {
			continue
		}
		edges = append(edges, [2]int{u, v})
	}
	prices := make([]int, m)
	for i := range prices {
		prices[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d %d %d %d\n", n, m, a, b, c)
	for i, v := range prices {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func runCase(bin string, input string) error {
	expect, err := runRef(input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runProg(bin, input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %s got %s", strings.TrimSpace(expect), strings.TrimSpace(got))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
