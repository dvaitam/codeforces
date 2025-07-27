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

func run(bin string, input string) (string, error) {
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

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveF(n int, tVals, xVals []int64) bool {
	for c := 0; c <= n; c++ {
		p := make([]int, 0, n+1)
		p = append(p, 0)
		for i := 1; i <= n; i++ {
			if i == c {
				continue
			}
			p = append(p, i)
		}
		ok := true
		for k := 0; k+1 < len(p); k++ {
			i := p[k]
			j := p[k+1]
			dt := tVals[j] - tVals[i]
			if c != 0 && i < c && c < j {
				need := abs(xVals[c]-xVals[i]) + abs(xVals[j]-xVals[c])
				if dt < need {
					ok = false
					break
				}
			} else {
				if dt < abs(xVals[j]-xVals[i]) {
					ok = false
					break
				}
			}
		}
		if !ok {
			continue
		}
		if c > 0 {
			prev := 0
			for i := 1; i < c; i++ {
				if i != c {
					prev = i
				}
			}
			if prev == p[len(p)-1] {
				if tVals[c]-tVals[prev] < abs(xVals[c]-xVals[prev]) {
					ok = false
				}
			}
		}
		if ok {
			return true
		}
	}
	return false
}

func runCase(bin string, n int, tVals, xVals []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tVals[i], xVals[i]))
	}
	expect := "NO"
	if solveF(n, tVals, xVals) {
		expect = "YES"
	}
	out, err := run(bin, sb.String())
	if err != nil {
		return err
	}
	if out != expect {
		return fmt.Errorf("expected %s got %s", expect, out)
	}
	return nil
}

func uniqueRandInts(rng *rand.Rand, n int, low, high int64) []int64 {
	m := make(map[int64]struct{})
	res := make([]int64, 0, n)
	for len(res) < n {
		v := rng.Int63n(high-low+1) + low
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		res = append(res, v)
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	total := 0
	// edge case
	tVals := []int64{0, 1}
	xVals := []int64{0, 0}
	if err := runCase(bin, 1, tVals, xVals); err != nil {
		fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
		os.Exit(1)
	}
	total++
	for total < 100 {
		n := rng.Intn(6) + 1
		tVals = make([]int64, n+1)
		xVals = make([]int64, n+1)
		curT := int64(0)
		for i := 1; i <= n; i++ {
			curT += int64(rng.Intn(5) + 1)
			tVals[i] = curT
		}
		coords := uniqueRandInts(rng, n, -10, 10)
		for i := 1; i <= n; i++ {
			xVals[i] = coords[i-1]
		}
		if err := runCase(bin, n, tVals, xVals); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", total+1, err)
			os.Exit(1)
		}
		total++
	}
	fmt.Printf("All %d tests passed\n", total)
}
