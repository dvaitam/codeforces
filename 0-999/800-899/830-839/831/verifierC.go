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

type testCaseC struct {
	k int
	n int
	a []int
	b []int
}

func generateCaseC(rng *rand.Rand) (string, testCaseC) {
	k := rng.Intn(10) + 1
	n := rng.Intn(k) + 1
	a := make([]int, k)
	for i := 0; i < k; i++ {
		a[i] = rng.Intn(401) - 200 // -200..200
	}
	// compute prefix sums
	pref := make([]int, k)
	sum := 0
	for i := 0; i < k; i++ {
		sum += a[i]
		pref[i] = sum
	}
	offset := rng.Intn(1000) - 500
	set := make(map[int]struct{})
	b := make([]int, 0, n)
	for len(b) < n {
		if rng.Float64() < 0.8 {
			val := pref[rng.Intn(k)] + offset
			if _, ok := set[val]; !ok {
				set[val] = struct{}{}
				b = append(b, val)
			}
		} else {
			val := rng.Intn(2000) - 1000
			if _, ok := set[val]; !ok {
				set[val] = struct{}{}
				b = append(b, val)
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", k, n)
	for i := 0; i < k; i++ {
		fmt.Fprintf(&sb, "%d ", a[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", b[i])
	}
	sb.WriteByte('\n')
	return sb.String(), testCaseC{k: k, n: n, a: a, b: b}
}

func expectedC(tc testCaseC) int {
	pref := make([]int, tc.k)
	sum := 0
	for i := 0; i < tc.k; i++ {
		sum += tc.a[i]
		pref[i] = sum
	}
	prefSet := make(map[int]struct{})
	for _, v := range pref {
		prefSet[v] = struct{}{}
	}
	candidate := make(map[int]struct{})
	for _, v := range pref {
		candidate[tc.b[0]-v] = struct{}{}
	}
	count := 0
	for x := range candidate {
		ok := true
		for _, bj := range tc.b {
			if _, ok2 := prefSet[bj-x]; !ok2 {
				ok = false
				break
			}
		}
		if ok {
			count++
		}
	}
	return count
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input, tc := generateCaseC(rng)
		expect := expectedC(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		var val int
		if _, err := fmt.Sscan(got, &val); err != nil || val != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
