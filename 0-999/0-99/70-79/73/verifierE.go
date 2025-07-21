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

func solve(n int, x int, a []int) int {
	if x <= 2 {
		return 0
	}
	hasOne := false
	for _, v := range a {
		if v == 1 {
			hasOne = true
			break
		}
	}
	if hasOne {
		return 1
	}
	L := x - 1
	const maxL = 2000000
	if L > maxL {
		return -1
	}
	uniq := make([]bool, L+1)
	for _, v := range a {
		if v <= L {
			uniq[v] = true
		}
	}
	vals := make([]int, 0)
	for v, ok := range uniq {
		if ok && v >= 2 {
			vals = append(vals, v)
		}
	}
	if len(vals) == 0 {
		return -1
	}
	sort.Ints(vals)
	spf := make([]int, L+1)
	for i := 2; i <= L; i++ {
		if spf[i] == 0 {
			for j := i; j <= L; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	hasB := make([]bool, L+1)
	B := make([]int, 0)
	for _, v := range vals {
		ok := true
		divs := []int{1}
		vv := v
		for vv > 1 {
			p := spf[vv]
			cnt := 0
			for vv%p == 0 {
				vv /= p
				cnt++
			}
			base := len(divs)
			mul := 1
			for e := 1; e <= cnt; e++ {
				mul *= p
				for i := 0; i < base; i++ {
					divs = append(divs, divs[i]*mul)
				}
			}
		}
		for _, d := range divs {
			if d == v {
				continue
			}
			if hasB[d] {
				ok = false
				break
			}
		}
		if ok {
			hasB[v] = true
			B = append(B, v)
		}
	}
	covered := make([]bool, L+1)
	for _, b := range B {
		for m := b; m <= L; m += b {
			covered[m] = true
		}
	}
	for k := 2; k <= L; k++ {
		if !covered[k] {
			return -1
		}
	}
	return len(B)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		x := rng.Intn(15) + 3
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(x) + 1
		}
		expected := fmt.Sprintf("%d", solve(n, x, a))
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for j, v := range a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
