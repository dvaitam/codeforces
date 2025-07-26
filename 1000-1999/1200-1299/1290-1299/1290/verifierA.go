package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expectedA(n, m, k int, a []int) int {
	if k > m-1 {
		k = m - 1
	}
	uncontrolled := m - 1 - k
	finalLen := n - m + 1
	ans := 0
	for takeFront := 0; takeFront <= k; takeFront++ {
		best := int(^uint(0) >> 1)
		for skipFront := 0; skipFront <= uncontrolled; skipFront++ {
			l := takeFront + skipFront
			r := l + finalLen - 1
			cand := a[l]
			if a[r] > cand {
				cand = a[r]
			}
			if cand < best {
				best = cand
			}
		}
		if best > ans {
			ans = best
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		n := rng.Intn(8) + 1
		m := rng.Intn(n) + 1
		k := rng.Intn(n)
		if k > n-1 {
			k = n - 1
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(10)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", a[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := expectedA(n, m, k, a)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", t, want, strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
