package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func prefixSum(k, n, D int64) int64 {
	if k < 0 {
		return 0
	}
	j1 := D - 1
	j3 := n - D
	u1 := minInt64(k, j1)
	sum1 := (u1 + 1) * (u1 + 1)
	u2 := minInt64(k, j3)
	var sum2 int64
	if u2 > j1 {
		L := j1 + 1
		R := u2
		cnt := R - L + 1
		sumJ := (R*(R+1)/2 - (L-1)*L/2)
		sum2 = sumJ + cnt*D
	}
	var sum3 int64
	if k > j3 {
		cnt3 := k - j3
		sum3 = cnt3 * n
	}
	return sum1 + sum2 + sum3
}

func widthAt(t, n, D int64) int64 {
	w1 := 2*t + 1
	w2 := t + D
	w := w1
	if w2 < w {
		w = w2
	}
	if n < w {
		w = n
	}
	return w
}

func countOn(t, n, x, y int64) int64 {
	D := minInt64(y, n-y+1)
	w0 := widthAt(t, n, D)
	up := t
	if up > x-1 {
		up = x - 1
	}
	down := t
	if down > n-x {
		down = n - x
	}
	S_t1 := prefixSum(t-1, n, D)
	S_up := prefixSum(t-up-1, n, D)
	S_down := prefixSum(t-down-1, n, D)
	sumRows := (S_t1 - S_up) + (S_t1 - S_down)
	return w0 + sumRows
}

func solve(n, x, y, c int64) int64 {
	var maxDist int64
	corners := []struct{ dx, dy int64 }{{x - 1, y - 1}, {x - 1, n - y}, {n - x, y - 1}, {n - x, n - y}}
	for _, d := range corners {
		s := d.dx + d.dy
		if s > maxDist {
			maxDist = s
		}
	}
	lo, hi := int64(-1), maxDist
	for lo+1 < hi {
		mid := (lo + hi) / 2
		if countOn(mid, n, x, y) >= c {
			hi = mid
		} else {
			lo = mid
		}
	}
	if hi < 0 {
		hi = 0
	}
	return hi
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) != 4 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.ParseInt(fields[0], 10, 64)
		xVal, _ := strconv.ParseInt(fields[1], 10, 64)
		yVal, _ := strconv.ParseInt(fields[2], 10, 64)
		cVal, _ := strconv.ParseInt(fields[3], 10, 64)
		input := fmt.Sprintf("%d %d %d %d\n", nVal, xVal, yVal, cVal)
		want := solve(nVal, xVal, yVal, cVal)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		gotVal, _ := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if gotVal != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, want, gotVal)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
