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

func build(seq *[][]int) bool {
	v := *seq
	a := v[len(v)-2]
	b := v[len(v)-1]
	newB := make([]int, len(b)+1)
	copy(newB[1:], b)
	mult := 1
	for i := 0; i < len(a); i++ {
		if mult == -1 && ((a[i] == -1 && newB[i] == 1) || (a[i] == 1 && newB[i] == -1)) {
			return false
		}
		if (a[i] == 1 && newB[i] == 1) || (a[i] == -1 && newB[i] == -1) {
			mult = -1
		}
	}
	for i := 0; i < len(a); i++ {
		newB[i] += mult * a[i]
	}
	*seq = append(*seq, newB)
	return true
}

func solveB(n int) (string, bool) {
	seq := make([][]int, 0, n+2)
	seq = append(seq, []int{0})
	seq = append(seq, []int{1})
	seq = append(seq, []int{0, 1})
	seq = append(seq, []int{1, 0, 1})
	for len(seq) < n+2 {
		if !build(&seq) {
			return "-1", false
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range seq[n+1] {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", n-1))
	for i, v := range seq[n] {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return strings.TrimSpace(sb.String()), true
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
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(10) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		input := sb.String()
		expected, ok := solveB(n)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		if !ok {
			if strings.TrimSpace(got) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %q\n", t+1, got)
				os.Exit(1)
			}
			continue
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", t+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
