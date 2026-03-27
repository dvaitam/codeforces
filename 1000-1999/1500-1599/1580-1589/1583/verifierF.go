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

type caseF struct{ n, k int }

func genCase(rng *rand.Rand) caseF {
	n := rng.Intn(8) + 2
	k := rng.Intn(n-1) + 2
	return caseF{n, k}
}

func solve1583F(tc caseF) []string {
	n, k := tc.n, tc.k

	c := 0
	temp := n - 1
	for temp > 0 {
		temp /= k
		c++
	}

	res := []string{strconv.Itoa(c)}

	var sb strings.Builder
	first := true
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			d := 1
			u, v := i, j
			for u/k != v/k {
				u /= k
				v /= k
				d++
			}
			if !first {
				sb.WriteByte(' ')
			}
			first = false
			sb.WriteString(strconv.Itoa(d))
		}
	}
	res = append(res, sb.String())
	return res
}

func runCase(bin string, tc caseF) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	exp := solve1583F(tc)
	need := 1 + tc.n*(tc.n-1)/2
	if len(fields) != need {
		return fmt.Errorf("expected %d numbers got %d", need, len(fields))
	}
	if fields[0] != exp[0] {
		return fmt.Errorf("expected c=%s got %s", exp[0], fields[0])
	}
	expColors := strings.Fields(exp[1])
	for i := 1; i < need; i++ {
		if fields[i] != expColors[i-1] {
			return fmt.Errorf("color %d expected %s got %s", i, expColors[i-1], fields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
