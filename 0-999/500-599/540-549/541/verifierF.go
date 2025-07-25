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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, T int, tasks []struct{ t, q int }) string {
	buckets := make([][]int, T)
	for _, task := range tasks {
		L := T - task.t
		if L < 0 {
			continue
		}
		buckets[L] = append(buckets[L], task.q)
	}
	prefs := make([][]int, T)
	for d := 0; d < T; d++ {
		b := buckets[d]
		sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
		pref := make([]int, len(b)+1)
		for i, v := range b {
			pref[i+1] = pref[i] + v
		}
		prefs[d] = pref
	}
	const neg = -1000000000
	prev := make([]int, n+1)
	for i := range prev {
		prev[i] = neg
	}
	prev[1] = 0
	for d := 0; d < T; d++ {
		next := make([]int, n+1)
		for i := range next {
			next[i] = neg
		}
		pref := prefs[d]
		m := len(pref) - 1
		for s := 0; s <= n; s++ {
			base := prev[s]
			if base < 0 {
				continue
			}
			maxTake := m
			if s < maxTake {
				maxTake = s
			}
			for x := 0; x <= maxTake; x++ {
				profit := base + pref[x]
				newSlots := (s - x) * 2
				if newSlots > n {
					newSlots = n
				}
				if profit > next[newSlots] {
					next[newSlots] = profit
				}
			}
		}
		prev = next
	}
	ans := 0
	for _, v := range prev {
		if v > ans {
			ans = v
		}
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		T := rng.Intn(10) + 1
		tasks := make([]struct{ t, q int }, n)
		for j := 0; j < n; j++ {
			t := rng.Intn(T) + 1
			q := rng.Intn(10) + 1
			tasks[j] = struct{ t, q int }{t, q}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, T))
		for _, task := range tasks {
			sb.WriteString(fmt.Sprintf("%d %d\n", task.t, task.q))
		}
		input := sb.String()
		exp := expected(n, T, tasks)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
