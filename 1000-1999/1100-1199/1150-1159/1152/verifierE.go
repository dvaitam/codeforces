package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkAnswer(ans []int, b, c []int) bool {
	if len(ans) != len(b)+1 {
		return false
	}
	n := len(ans)
	pairs1 := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		x, y := ans[i], ans[i+1]
		if x <= 0 || y <= 0 || x > 1_000_000_000 || y > 1_000_000_000 {
			return false
		}
		if x > y {
			x, y = y, x
		}
		pairs1[i] = [2]int{x, y}
	}
	pairs2 := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		x, y := b[i], c[i]
		if x > y {
			x, y = y, x
		}
		pairs2[i] = [2]int{x, y}
	}
	sort.Slice(pairs1, func(i, j int) bool {
		if pairs1[i][0] == pairs1[j][0] {
			return pairs1[i][1] < pairs1[j][1]
		}
		return pairs1[i][0] < pairs1[j][0]
	})
	sort.Slice(pairs2, func(i, j int) bool {
		if pairs2[i][0] == pairs2[j][0] {
			return pairs2[i][1] < pairs2[j][1]
		}
		return pairs2[i][0] < pairs2[j][0]
	})
	for i := 0; i < n-1; i++ {
		if pairs1[i] != pairs2[i] {
			return false
		}
	}
	return true
}

func generateCase(rng *rand.Rand) (string, []int, []int) {
	n := rng.Intn(8) + 2
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(1000) + 1
	}
	b := make([]int, n-1)
	c := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		if a[i] < a[i+1] {
			b[i] = a[i]
			c[i] = a[i+1]
		} else {
			b[i] = a[i+1]
			c[i] = a[i]
		}
	}
	perm := rng.Perm(n - 1)
	b2 := make([]int, n-1)
	c2 := make([]int, n-1)
	for i, p := range perm {
		b2[i] = b[p]
		c2[i] = c[p]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range b2 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range c2 {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), b2, c2
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// deterministic cases
	fixed := []struct {
		b []int
		c []int
	}{
		{[]int{3}, []int{4}},
		{[]int{2, 3}, []int{5, 4}},
	}
	idx := 0
	for ; idx < len(fixed); idx++ {
		n := len(fixed[idx].b) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range fixed[idx].b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i, v := range fixed[idx].c {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		toks := strings.Fields(out)
		if len(toks) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", idx+1, n, len(toks))
			os.Exit(1)
		}
		ans := make([]int, n)
		for i := 0; i < n; i++ {
			val, err := strconv.Atoi(toks[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid number %v\n", idx+1, err)
				os.Exit(1)
			}
			ans[i] = val
		}
		if !checkAnswer(ans, fixed[idx].b, fixed[idx].c) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid answer\n", idx+1)
			os.Exit(1)
		}
	}
	for ; idx < 100; idx++ {
		inp, b, c := generateCase(rng)
		out, err := run(bin, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, inp)
			os.Exit(1)
		}
		toks := strings.Fields(out)
		n := len(b) + 1
		if len(toks) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", idx+1, n, len(toks))
			os.Exit(1)
		}
		ans := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(toks[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid number %v\n", idx+1, err)
				os.Exit(1)
			}
			ans[i] = v
		}
		if !checkAnswer(ans, b, c) {
			fmt.Fprintf(os.Stderr, "case %d failed: wrong answer\ninput:%soutput:%s\n", idx+1, inp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
