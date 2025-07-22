package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveE(r *bufio.Reader) string {
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	vals := make([]int, 0)
	seen := make(map[int]bool)
	for _, v := range a {
		if !seen[v] {
			vals = append(vals, v)
			seen[v] = true
		}
	}
	sort.Ints(vals)
	bestLen := 0
	bestX, bestY := 0, 0
	for i := 0; i < len(vals); i++ {
		x := vals[i]
		count := 0
		for _, v := range a {
			if v == x {
				count++
			}
		}
		l := (count / 4) * 4
		if l > bestLen {
			bestLen = l
			bestX = x
			bestY = x
		}
	}
	for i := 0; i < len(vals); i++ {
		for j := 0; j < len(vals); j++ {
			if i == j {
				continue
			}
			x := vals[i]
			y := vals[j]
			state := 0
			cnt := 0
			for _, v := range a {
				if state == 0 || state == 2 {
					if v == x {
						state++
					}
				} else {
					if v == y {
						state++
						if state == 4 {
							cnt += 4
							state = 0
						}
					}
				}
			}
			if cnt > bestLen {
				bestLen = cnt
				bestX = x
				bestY = y
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", bestLen))
	for i := 0; i < bestLen/4; i++ {
		fmt.Fprintf(&sb, "%d %d %d %d", bestX, bestY, bestX, bestY)
		if i+1 < bestLen/4 {
			sb.WriteByte(' ')
		}
	}
	if bestLen > 0 {
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", rng.Intn(10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		expect := solveE(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
