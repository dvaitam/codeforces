package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `2 2 1 11 11 16 12 34 18
3 1 2 0 0 1 18 20 22 29 4 3 14 18
3 1 2 1 1 0 18 16 32 32 9 2 26 2
1 3 2 011 8 11 10 17 19 8 26 12
3 2 1 01 10 11 4 18 14 35
1 3 3 110 20 13 30 31 8 10 13 16 6 2 25 10
2 1 1 0 0 2 3 19 15
3 3 2 001 111 101 20 4 35 22 11 7 18 7
3 2 1 01 01 10 4 5 11 6
3 3 3 000 001 010 2 20 2 26 6 4 21 10 2 1 19 14
3 1 2 0 0 0 10 12 23 17 2 17 16 18
3 1 3 1 0 1 12 16 30 21 7 2 12 7 11 17 19 20
3 2 3 00 11 11 13 9 17 26 1 15 3 25 2 18 10 22
1 2 2 11 19 20 23 29 13 14 33 16
1 3 1 100 8 15 20 33
2 1 2 1 0 6 15 8 23 6 15 22 30
3 3 1 011 110 100 5 1 17 14
2 1 1 0 0 17 20 20 26
1 3 3 011 6 4 21 16 3 1 11 15 4 9 8 29
3 3 3 100 100 001 18 11 29 29 2 20 22 35 15 14 26 31
1 1 2 1 1 5 5 13 11 11 22 13
2 3 1 010 011 13 18 17 27
1 2 3 00 10 6 26 8 10 13 20 22 14 4 17 21
2 2 2 10 10 16 14 17 23 11 5 16 25
3 2 3 00 00 00 13 1 16 13 18 17 27 31 16 19 22 32
1 2 1 10 14 7 25 10
1 3 1 100 16 13 24 19
3 1 1 0 0 0 15 13 26 30
1 1 3 1 5 19 17 39 14 17 29 27 16 16 36 22
3 3 1 011 100 101 19 10 34 12
1 3 1 000 2 10 2 24
2 1 1 1 1 17 13 33 29
1 3 1 010 10 18 29 31
2 2 3 00 00 10 17 28 25 11 3 26 11 10 14 22 26
1 1 3 0 8 10 18 11 2 16 15 20 16 20 18 24
2 2 1 11 10 4 16 8 16
1 3 3 010 18 12 24 24 16 4 17 23 15 20 35 30
3 1 3 1 0 1 10 4 26 10 2 13 16 24 7 15 18 35
1 1 1 1 9 1 25 19
3 1 1 0 1 1 4 5 17 23
2 1 1 1 0 4 14 8 14
2 2 3 10 11 9 3 20 5 4 12 4 23 12 6 12 13
2 1 3 0 0 1 7 4 7 10 12 10 31 8 5 13 19
1 2 2 10 1 7 12 17 16 10 25 27
3 2 1 00 10 10 5 8 15 24
1 1 1 1 12 14 13 18
3 1 2 0 0 0 14 10 31 23 5 19 18 28
3 2 1 01 10 11 1 14 11 28
1 2 2 10 6 4 14 7 18 20 22 34
2 1 2 1 0 8 15 18 31 5 12 19 32
3 1 2 0 1 0 15 20 29 20 7 10 10 30
2 3 3 011 010 18 13 26 33 1 4 9 5 1 9 13 25
3 3 2 101 110 000 9 10 26 20 4 17 11 22
1 2 2 10 19 17 39 23 18 4 31 24
3 2 3 11 11 00 4 4 16 16 19 15 23 32 10 12 30 27
3 2 1 11 11 01 10 5 25 6
3 1 1 1 1 1 1 17 3 19
3 3 3 101 000 110 5 19 14 25 4 14 18 24 13 6 23 19
3 3 2 010 100 011 13 14 28 26 8 7 22 13
3 3 1 100 001 000 20 10 39 12
3 3 2 111 011 111 7 11 15 12 2 2 7 13
1 2 3 00 3 14 10 33 13 18 20 32 7 11 26 14
3 1 2 1 1 1 9 1 25 2 7 12 9 18
3 2 1 01 11 11 16 12 23 13
2 3 1 011 101 16 15 30 18
1 1 1 0 5 14 11 28
3 1 2 1 0 0 8 16 15 20 9 12 19 25
1 3 2 011 17 20 31 37 9 9 16 9
1 3 3 001 8 7 17 7 18 17 31 18 4 13 24 21
1 3 3 101 8 8 10 24 10 11 17 22 16 10 34 15
1 1 3 1 12 19 32 19 5 13 9 18 17 3 21 9
2 3 3 000 011 20 19 24 39 16 4 35 4 17 20 28 35
2 2 1 00 11 18 11 20 19
1 3 2 011 13 2 19 3 11 8 21 22
3 3 3 011 010 100 12 1 27 21 2 1 9 2 1 8 21 18
1 1 2 1 5 7 19 20 5 12 14 17
3 2 3 11 01 11 2 19 5 32 13 6 13 22 5 20 21 24
1 2 1 00 1 6 18 11
3 1 2 0 1 0 20 20 21 28 11 13 11 33
1 2 1 11 5 15 12 31
2 1 3 1 1 9 16 21 16 10 17 19 34 16 2 33 20
3 2 3 01 10 11 16 2 16 10 2 9 20 18 7 17 23 27
2 2 1 01 01 6 5 16 5
3 1 3 0 1 1 10 10 20 25 13 20 26 25 1 5 19 6
2 1 2 0 1 9 20 15 22 18 14 26 19
3 1 1 0 0 1 14 9 23 18
1 2 3 11 18 17 35 27 11 7 24 11 1 17 5 35
2 2 2 01 00 12 17 17 23 12 16 12 23
3 1 2 0 1 0 19 15 26 29 17 4 23 9
2 1 2 1 1 9 14 20 33 11 3 20 3
2 1 2 0 1 13 14 33 34 13 2 31 16
2 3 1 110 110 2 3 20 14
2 1 1 0 0 18 16 19 26
1 2 2 01 14 5 33 9 13 10 29 11
1 1 1 1 2 17 3 34
3 3 3 101 000 101 11 9 19 25 15 5 29 22 5 2 25 20
1 3 3 010 7 15 26 22 15 17 20 27 5 16 22 17
3 1 3 1 0 0 4 14 23 25 19 15 29 27 17 12 37 15
1 2 1 00 1 11 20 17
1 2 3 01 13 2 32 7 12 3 25 4 15 12 34 31`

// referenceSolve replicates 1186E.go.
func referenceSolve(n, m, q int, rows []string, queries [][4]int64) []int64 {
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			row[j] = int(rows[i][j] - '0')
		}
		a[i] = row
	}
	P := make([][]int64, n+1)
	for i := range P {
		P[i] = make([]int64, m+1)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			P[i+1][j+1] = P[i+1][j] + P[i][j+1] - P[i][j] + int64(a[i][j])
		}
	}
	total := int64(n * m)

	var tmOnes func(int64) int64
	tmOnes = func(k int64) int64 {
		if k <= 0 {
			return 0
		}
		b := 63 - bits.LeadingZeros64(uint64(k))
		p := int64(1) << b
		if p == k {
			return p / 2
		}
		r := k - p
		return p/2 + (r - tmOnes(r))
	}

	sumRect := func(x, y int64) int64 {
		if x <= 0 || y <= 0 {
			return 0
		}
		U := x / int64(n)
		R := x % int64(n)
		V := y / int64(m)
		C := y % int64(m)
		U1 := tmOnes(U)
		U0 := U - U1
		V1 := tmOnes(V)
		V0 := V - V1
		countEven := U0*V0 + U1*V1
		countOdd := U0*V1 + U1*V0
		res := countEven*P[n][m] + countOdd*(total-P[n][m])
		if C > 0 {
			sumCols := P[n][C]
			blockSize := int64(n) * C
			pv := bits.OnesCount64(uint64(V)) & 1
			var uSame int64
			if pv == 0 {
				uSame = U0
			} else {
				uSame = U1
			}
			uDiff := U - uSame
			res += uSame*sumCols + uDiff*(blockSize-sumCols)
		}
		if R > 0 {
			sumRows := P[R][m]
			blockSize := R * int64(m)
			pu := bits.OnesCount64(uint64(U)) & 1
			var vSame int64
			if pu == 0 {
				vSame = V0
			} else {
				vSame = V1
			}
			vDiff := V - vSame
			res += vSame*sumRows + vDiff*(blockSize-sumRows)
		}
		if R > 0 && C > 0 {
			sumCorner := P[R][C]
			blockSize := R * C
			pu := bits.OnesCount64(uint64(U)) & 1
			pv := bits.OnesCount64(uint64(V)) & 1
			if pu == pv {
				res += sumCorner
			} else {
				res += blockSize - sumCorner
			}
		}
		return res
	}

	results := make([]int64, q)
	for i, qu := range queries {
		x1, y1, x2, y2 := qu[0], qu[1], qu[2], qu[3]
		results[i] = sumRect(x2, y2) - sumRect(x1-1, y2) - sumRect(x2, y1-1) + sumRect(x1-1, y1-1)
	}
	return results
}

func parseLine(line string) (int, int, int, []string, [][4]int64, error) {
	fields := strings.Fields(line)
	if len(fields) < 3 {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid n: %v", err)
	}
	m, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid m: %v", err)
	}
	q, err := strconv.Atoi(fields[2])
	if err != nil {
		return 0, 0, 0, nil, nil, fmt.Errorf("invalid q: %v", err)
	}
	if len(fields) != 3+n+4*q {
		return 0, 0, 0, nil, nil, fmt.Errorf("expected %d fields got %d", 3+n+4*q, len(fields))
	}
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		rows[i] = fields[3+i]
		if len(rows[i]) != m {
			return 0, 0, 0, nil, nil, fmt.Errorf("row %d length mismatch", i+1)
		}
	}
	queries := make([][4]int64, q)
	for i := 0; i < q; i++ {
		a, _ := strconv.ParseInt(fields[3+n+4*i], 10, 64)
		b, _ := strconv.ParseInt(fields[3+n+4*i+1], 10, 64)
		c, _ := strconv.ParseInt(fields[3+n+4*i+2], 10, 64)
		d, _ := strconv.ParseInt(fields[3+n+4*i+3], 10, 64)
		queries[i] = [4]int64{a, b, c, d}
	}
	return n, m, q, rows, queries, nil
}

func runCase(bin string, line string) error {
	n, m, q, rows, queries, err := parseLine(line)
	if err != nil {
		return err
	}
	expectVals := referenceSolve(n, m, q, rows, queries)
	expectFields := make([]string, len(expectVals))
	for i, v := range expectVals {
		expectFields[i] = strconv.FormatInt(v, 10)
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d %d\n", n, m, q)
	for _, r := range rows {
		input.WriteString(r)
		input.WriteByte('\n')
	}
	for _, qu := range queries {
		fmt.Fprintf(&input, "%d %d %d %d\n", qu[0], qu[1], qu[2], qu[3])
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotFields := strings.Fields(strings.TrimSpace(out.String()))
	if len(gotFields) != len(expectFields) {
		return fmt.Errorf("expected %d outputs got %d", len(expectFields), len(gotFields))
	}
	for i := range expectFields {
		if gotFields[i] != expectFields[i] {
			return fmt.Errorf("query %d: expected %s got %s", i+1, expectFields[i], gotFields[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
