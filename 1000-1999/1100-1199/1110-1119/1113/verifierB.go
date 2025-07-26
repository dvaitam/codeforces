package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func expected(arr []int) int {
	var b [101]int
	ans := 0
	for _, v := range arr {
		if v >= 1 && v <= 100 {
			b[v]++
		}
		ans += v
	}
	adj := make([][]int, 101)
	for j := 1; j <= 100; j++ {
		for d := 2; d*d <= j; d++ {
			if j%d == 0 {
				adj[j] = append(adj[j], d)
				if d*d != j {
					adj[j] = append(adj[j], j/d)
				}
			}
		}
	}
	best := ans
	for j := 1; j <= 100; j++ {
		if b[j] == 0 || len(adj[j]) == 0 {
			continue
		}
		for _, d := range adj[j] {
			for i := 1; i <= 100; i++ {
				if b[i] == 0 {
					continue
				}
				pre := i + j
				cur := i*d + j/d
				if cur < pre {
					cand := ans - pre + cur
					if cand < best {
						best = cand
					}
				}
			}
		}
	}
	return best
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	file, err := os.Open(filepath.Join(dir, "testcasesB.txt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read test count: %v\n", err)
		os.Exit(1)
	}
	for caseNum := 1; caseNum <= t; caseNum++ {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read n on case %d: %v\n", caseNum, err)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i])
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", arr[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := fmt.Sprintf("%d", expected(arr))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", caseNum, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
