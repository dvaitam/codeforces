package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
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

func solveCase(sets []uint64) string {
	var all uint64
	for _, m := range sets {
		all |= m
	}
	ans := 0
	for e := 0; e < 50; e++ {
		if all&(1<<uint(e)) == 0 {
			continue
		}
		var u uint64
		for _, m := range sets {
			if m&(1<<uint(e)) == 0 {
				u |= m
			}
		}
		cnt := bits.OnesCount64(u)
		if cnt > ans {
			ans = cnt
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	sets := make([]uint64, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		k := rng.Intn(6) + 1
		var mask uint64
		elems := make([]int, 0, k)
		for len(elems) < k {
			x := rng.Intn(50) + 1
			dup := false
			for _, v := range elems {
				if v == x {
					dup = true
					break
				}
			}
			if dup {
				continue
			}
			elems = append(elems, x)
			mask |= 1 << uint(x-1)
		}
		sets[i] = mask
		sb.WriteString(fmt.Sprintf("%d", k))
		for _, v := range elems {
			sb.WriteString(fmt.Sprintf(" %d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), solveCase(sets)
}

func fixedCases() [][2]string {
	return [][2]string{
		{"1\n1\n1 1\n", "0"},
		{"1\n2\n1 1\n1 2\n", "1"},
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range fixedCases() {
		out, err := runCandidate(bin, tc[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: %v\ninput:\n%s", idx+1, err, tc[0])
			os.Exit(1)
		}
		if out != tc[1] {
			fmt.Fprintf(os.Stderr, "fixed case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc[1], out, tc[0])
			os.Exit(1)
		}
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
