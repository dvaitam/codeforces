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

func solve(arr []int) string {
	b := append([]int(nil), arr...)
	sort.Ints(b)
	m := b[0]
	for i := range arr {
		if arr[i] != b[i] && arr[i]%m != 0 {
			return "NO"
		}
	}
	return "YES"
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		sb.WriteString(fmt.Sprint(v))
		if i+1 < len(arr) {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	got, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	got = strings.ToUpper(strings.TrimSpace(got))
	exp := solve(arr)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
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
		n := rng.Intn(30) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(1_000_000) + 1
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
