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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(arr []int64) string {
	if len(arr) == 0 {
		return "-1"
	}
	g := arr[0]
	for _, v := range arr[1:] {
		g = gcd(g, v)
	}
	if g != arr[0] {
		return "-1"
	}
	var sb strings.Builder
	n := int64(len(arr))
	fmt.Fprintf(&sb, "%d\n", 2*n-1)
	sb.WriteString(fmt.Sprintf("%d", arr[0]))
	for i := 1; i < len(arr); i++ {
		sb.WriteString(fmt.Sprintf(" %d %d", arr[i], arr[0]))
	}
	sb.WriteByte('\n')
	return strings.TrimSpace(sb.String())
}

func genValidCase(rng *rand.Rand) []int64 {
	m := rng.Intn(5) + 1
	base := int64(rng.Intn(1000) + 1)
	set := make(map[int64]struct{})
	arr := make([]int64, 0, m)
	arr = append(arr, base)
	set[base] = struct{}{}
	for len(arr) < m {
		v := base * int64(rng.Intn(20)+1)
		if _, ok := set[v]; ok {
			continue
		}
		set[v] = struct{}{}
		arr = append(arr, v)
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr
}

func genInvalidCase(rng *rand.Rand) []int64 {
	arr := genValidCase(rng)
	if len(arr) == 1 {
		arr[0]++
		return arr
	}
	arr[0]++
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return arr
}

func buildInput(arr []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(arr))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		var arr []int64
		if t%2 == 0 {
			arr = genValidCase(rng)
		} else {
			arr = genInvalidCase(rng)
		}
		input := buildInput(arr)
		expected := solve(arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
