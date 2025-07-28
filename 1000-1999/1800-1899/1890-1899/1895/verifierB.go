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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveB(n int, arr []int) string {
	sort.Ints(arr)
	pairs := make([][2]int, n)
	total := n * 2
	for i := 0; i < n; i++ {
		pairs[i][0] = arr[i]
		pairs[i][1] = arr[total-1-i]
	}
	ans := 0
	for i := 0; i+1 < n; i++ {
		ans += abs(pairs[i][0] - pairs[i+1][0])
		ans += abs(pairs[i][1] - pairs[i+1][1])
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", ans))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", pairs[i][0], pairs[i][1]))
	}
	return strings.TrimSpace(sb.String())
}

func genCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(4) + 1
	arr := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		arr[i] = rng.Intn(100)
	}
	return n, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr := genCase(rng)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		expect := solveB(n, append([]int(nil), arr...))
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\n got:\n%s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
