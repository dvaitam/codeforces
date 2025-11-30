package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []string{
	"14 9 9 15",
	"4 3 4 16",
	"41 38 15 245",
	"55 27 43 343",
	"76 33 70 2817",
	"93 20 6 152",
	"17 13 12 15",
	"48 22 34 1635",
	"70 15 14 4537",
	"25 7 11 569",
	"52 10 3 341",
	"5 3 2 14",
	"73 15 53 1698",
	"44 6 7 1913",
	"93 18 35 3861",
	"85 22 28 6949",
	"50 28 13 833",
	"2 2 2 2",
	"66 14 47 2728",
	"33 33 15 1022",
	"75 28 73 420",
	"73 14 34 3541",
	"83 13 5 819",
	"42 22 8 1273",
	"37 7 16 621",
	"74 9 13 3208",
	"62 10 59 1330",
	"29 5 8 142",
	"6 6 5 4",
	"60 31 36 2261",
	"31 14 28 254",
	"28 20 27 100",
	"43 12 19 209",
	"58 29 13 2660",
	"61 17 51 1962",
	"25 15 22 164",
	"78 28 73 3310",
	"45 28 17 800",
	"22 5 20 278",
	"63 8 52 961",
	"78 45 73 2266",
	"1 1 1 1",
	"51 15 19 2373",
	"41 2 24 1670",
	"5 4 1 19",
	"81 71 3 6138",
	"7 1 4 40",
	"39 29 22 1319",
	"4 4 1 14",
	"36 3 22 92",
	"8 6 2 2",
	"81 68 81 805",
	"74 70 5 1739",
	"76 75 53 2291",
	"13 2 3 32",
	"20 18 8 119",
	"62 59 39 3744",
	"34 26 14 105",
	"67 63 10 1872",
	"17 13 15 34",
	"32 15 8 319",
	"56 56 44 2011",
	"54 17 39 787",
	"71 5 71 3199",
	"20 16 19 72",
	"7 3 4 38",
	"62 43 30 2443",
	"96 54 43 8733",
	"30 27 3 307",
	"69 65 31 2786",
	"93 39 31 8377",
	"56 18 19 1559",
	"42 21 14 1727",
	"13 5 5 50",
	"14 4 7 137",
	"21 11 7 215",
	"60 56 55 1003",
	"36 34 24 804",
	"9 8 8 13",
	"22 21 4 414",
	"68 35 34 71",
	"50 26 13 102",
	"78 21 8 3408",
	"55 22 28 2956",
	"49 3 18 1940",
	"85 74 76 117",
	"12 1 3 144",
	"31 6 29 881",
	"93 43 55 2290",
	"65 38 63 2022",
	"82 10 35 1205",
	"80 25 32 1102",
	"2 1 2 3",
	"39 34 2 781",
	"60 52 40 3158",
	"84 22 5 1370",
	"84 2 67 3273",
	"46 6 21 107",
	"87 60 66 1902",
	"34 25 4 887",
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

func referenceSolve(n, x, y, c int64) string {
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
	return strconv.FormatInt(hi, 10)
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func parseCase(line string) (int64, int64, int64, int64, error) {
	parts := strings.Fields(line)
	if len(parts) != 4 {
		return 0, 0, 0, 0, fmt.Errorf("expected 4 numbers, got %d", len(parts))
	}
	var vals [4]int64
	for i := 0; i < 4; i++ {
		v, err := strconv.ParseInt(parts[i], 10, 64)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("parse value %d: %w", i+1, err)
		}
		vals[i] = v
	}
	return vals[0], vals[1], vals[2], vals[3], nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	idx := 0
	for _, tc := range testcases {
		line := strings.TrimSpace(tc)
		if line == "" {
			continue
		}
		idx++
		n, x, y, c, err := parseCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d %d %d %d\n", n, x, y, c)
		expected := referenceSolve(n, x, y, c)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\ngot: %s\n", idx, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
