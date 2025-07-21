package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// solveCase computes whether a valid assignment exists and returns one if so.
func solveCase(n, m int, xs, ys, rods []int) (bool, []int) {
	x := make([]int, n+3)
	y := make([]int, n+3)
	for i := 1; i <= n; i++ {
		x[i] = xs[i-1]
		y[i] = ys[i-1]
	}
	x[0], y[0] = x[n], y[n]
	x[n+1], y[n+1] = x[1], y[1]
	x[n+2], y[n+2] = x[2], y[2]
	orig := make(map[int][]int)
	for i := 0; i < m; i++ {
		orig[rods[i]] = append(orig[rods[i]], i+1)
	}
	for start := 1; start <= 2; start++ {
		mp := make(map[int][]int, len(orig))
		for k, v := range orig {
			cp := make([]int, len(v))
			copy(cp, v)
			mp[k] = cp
		}
		ans := make([]int, n+1)
		for i := 1; i <= n; i++ {
			ans[i] = -1
		}
		bad := false
		for j := start; j <= n; j += 2 {
			t := abs(x[j]-x[j-1]) + abs(y[j]-y[j-1]) + abs(x[j]-x[j+1]) + abs(y[j]-y[j+1])
			list := mp[t]
			if len(list) == 0 {
				bad = true
				break
			}
			ans[j] = list[len(list)-1]
			mp[t] = list[:len(list)-1]
		}
		if !bad {
			return true, ans[1:]
		}
	}
	return false, nil
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func parseInput(input string) (n, m int, xs, ys, rods []int, err error) {
	rdr := bufio.NewReader(strings.NewReader(input))
	if _, err = fmt.Fscan(rdr, &n, &m); err != nil {
		return
	}
	xs = make([]int, n)
	ys = make([]int, n)
	for i := 0; i < n; i++ {
		if _, err = fmt.Fscan(rdr, &xs[i], &ys[i]); err != nil {
			return
		}
	}
	rods = make([]int, m)
	for i := 0; i < m; i++ {
		if _, err = fmt.Fscan(rdr, &rods[i]); err != nil {
			return
		}
	}
	return
}

func verify(input, output string) error {
	n, m, xs, ys, rods, err := parseInput(input)
	if err != nil {
		return fmt.Errorf("invalid input: %v", err)
	}
	expectedOK, _ := solveCase(n, m, xs, ys, rods)
	outLines := strings.Split(strings.TrimSpace(output), "\n")
	if len(outLines) == 0 {
		return fmt.Errorf("empty output")
	}
	first := strings.ToUpper(strings.TrimSpace(outLines[0]))
	if first != "YES" && first != "NO" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if first == "NO" {
		if expectedOK {
			return fmt.Errorf("expected YES got NO")
		}
		return nil
	}
	if !expectedOK {
		return fmt.Errorf("expected NO got YES")
	}
	if len(outLines) < 2 {
		return fmt.Errorf("missing second line")
	}
	fields := strings.Fields(outLines[1])
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	xExt := make([]int, n+3)
	yExt := make([]int, n+3)
	for j := 1; j <= n; j++ {
		xExt[j] = xs[j-1]
		yExt[j] = ys[j-1]
	}
	xExt[0], yExt[0] = xExt[n], yExt[n]
	xExt[n+1], yExt[n+1] = xExt[1], yExt[1]
	xExt[n+2], yExt[n+2] = xExt[2], yExt[2]
	used := make(map[int]bool)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("invalid number %q", fields[i])
		}
		if v != -1 {
			if v < 1 || v > m {
				return fmt.Errorf("rod index out of range")
			}
			if used[v] {
				return fmt.Errorf("rod %d used multiple times", v)
			}
			used[v] = true
			j := i + 1
			t := abs(xExt[j]-xExt[j-1]) + abs(yExt[j]-yExt[j-1]) + abs(xExt[j]-xExt[j+1]) + abs(yExt[j]-yExt[j+1])
			if rods[v-1] != t {
				return fmt.Errorf("nail %d uses wrong rod length", j)
			}
		}
	}
	return nil
}

func runCase(bin, tc string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verify(tc, out.String())
}

func generateCase(rng *rand.Rand) string {
	n := 4
	if rng.Intn(2) == 1 {
		n = 6
	}
	if n%2 == 1 {
		n++
	}
	w := rng.Intn(10) + 1
	h := rng.Intn(10) + 1
	xs := []int{0, w, w, 0}
	ys := []int{0, 0, h, h}
	if n == 6 {
		xs = append(xs, -w)
		ys = append(ys, h)
	}
	m := rng.Intn(5) + 1
	rods := make([]int, m)
	for i := 0; i < m; i++ {
		rods[i] = rng.Intn(4*(w+h)) + 1
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", len(xs), m)
	for i := 0; i < len(xs); i++ {
		fmt.Fprintf(&b, "%d %d\n", xs[i], ys[i])
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rods[i])
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
