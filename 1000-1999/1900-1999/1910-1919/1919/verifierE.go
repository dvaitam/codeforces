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

const mod = 998244353

func bruteCount(p []int) int {
	n := len(p)
	cnt := 0
	arr := make([]int, n)
	cmp := make([]int, n)
	for mask := 0; mask < (1 << n); mask++ {
		sum := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				sum += 1
				arr[i] = sum
			} else {
				sum -= 1
				arr[i] = sum
			}
		}
		copy(cmp, arr)
		sort.Ints(cmp)
		ok := true
		for i := 0; i < n; i++ {
			if cmp[i] != p[i] {
				ok = false
				break
			}
		}
		if ok {
			cnt++
		}
	}
	return cnt % mod
}

func runExe(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	var out strings.Builder
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		a := make([]int, n)
		sum := 0
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				a[j] = 1
			} else {
				a[j] = -1
			}
			sum += a[j]
		}
		pref := make([]int, n)
		s := 0
		for j := 0; j < n; j++ {
			s += a[j]
			pref[j] = s
		}
		sort.Ints(pref)
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			fmt.Fprintf(&sb, "%d ", pref[j])
		}
		sb.WriteByte('\n')
		count := bruteCount(pref)
		out.WriteString(fmt.Sprintf("%d\n", count))
	}
	return sb.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := genCase(rng)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
