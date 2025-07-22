package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type exam struct {
	grade int64
	cost  int64
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n int, r, avg int64, exams []exam) string {
	sort.Slice(exams, func(i, j int) bool { return exams[i].cost < exams[j].cost })
	total := int64(0)
	for _, e := range exams {
		total += e.grade
	}
	need := int64(n)*avg - total
	if need <= 0 {
		return "0"
	}
	essays := int64(0)
	for i := 0; i < n && need > 0; i++ {
		canAdd := r - exams[i].grade
		if canAdd <= 0 {
			continue
		}
		add := canAdd
		if need < add {
			add = need
		}
		need -= add
		essays += add * exams[i].cost
	}
	return fmt.Sprintf("%d", essays)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	r := int64(rng.Intn(10) + 1)
	avg := int64(rng.Intn(int(r)) + 1)
	exams := make([]exam, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, r, avg)
	for i := 0; i < n; i++ {
		exams[i].grade = int64(rng.Intn(int(r) + 1))
		exams[i].cost = int64(rng.Intn(10) + 1)
		fmt.Fprintf(&sb, "%d %d\n", exams[i].grade, exams[i].cost)
	}
	expect := solveCase(n, r, avg, exams)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
