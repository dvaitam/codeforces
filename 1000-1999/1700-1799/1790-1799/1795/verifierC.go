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

func expected(a, b []int64) string {
	n := len(a)
	pref := make([]int64, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			pref[i] = b[i]
		} else {
			pref[i] = pref[i-1] + b[i]
		}
	}
	diff := make([]int64, n+1)
	ans := make([]int64, n)
	for i := 0; i < n; i++ {
		x := a[i]
		if i > 0 {
			x += pref[i-1]
		}
		pos := sort.Search(n, func(j int) bool { return pref[j] > x })
		diff[i]++
		diff[pos]--
		if pos < n {
			var prev int64
			if pos > 0 {
				prev = pref[pos-1]
			}
			ans[pos] += x - prev
		}
	}
	curr := int64(0)
	for i := 0; i < n; i++ {
		curr += diff[i]
		ans[i] += curr * b[i]
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", ans[i]))
	}
	return sb.String()
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(10) + 1
		a := make([]int64, n)
		b := make([]int64, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			a[i] = int64(rng.Intn(20) + 1)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < n; i++ {
			b[i] = int64(rng.Intn(20) + 1)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", b[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expectedOut := expected(a, b)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tc+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
