package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) string {
	m := rng.Intn(6) + 5 // 5..10
	k := rng.Intn(6) + 7 // 7..12
	x := make([]int, m)
	v := make([]int, m)
	for i := 0; i < m; i++ {
		x[i] = rng.Intn(100) + 1
		v[i] = rng.Intn(20) + 1
	}
	positions := make([][]int, k)
	for t := 0; t < k; t++ {
		coords := make([]int, m)
		for i := 0; i < m; i++ {
			coords[i] = x[i] + t*v[i]
		}
		perm := rng.Perm(m)
		row := make([]int, m)
		for j, idx := range perm {
			row[j] = coords[idx]
		}
		positions[t] = row
	}
	y := rng.Intn(k-2) + 1
	idx := rng.Intn(m)
	orig := positions[y][idx]
	c := orig
	for c == orig {
		c = rng.Intn(200) + 1
	}
	positions[y][idx] = c

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", m, k))
	for t := 0; t < k; t++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", positions[t][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1545D.go")
	return runProg(ref, input)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		expect, err := runRef(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
