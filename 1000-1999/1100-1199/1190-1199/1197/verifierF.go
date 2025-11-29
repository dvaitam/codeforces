package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const mod = 998244353

const testcasesRaw = `100
3
3 5 1
3
1 3 1
1 1 2
2 2 2
0 0 0
0 1 1
0 1 0
1
2
2
1 1 1
1 2 1
0 0 0
1 1 0
1 1 0
2
3 1
2
2 1 3
1 1 2
0 1 1
0 1 1
1 0 1
1
2
2
1 1 1
1 2 3
1 1 1
0 1 1
0 1 0
2
3 5
2
1 2 2
2 3 1
0 0 1
1 0 0
0 1 0
2
5 3
2
2 3 2
1 1 3
0 1 1
0 0 0
1 0 0
4
3 2 1 4
2
2 2 3
4 3 2
1 1 1
1 0 0
1 0 1
3
2 1 4
2
3 3 1
2 1 2
1 1 1
1 1 1
0 0 1
3
3 3 4
2
3 3 2
2 3 2
1 1 0
1 1 1
0 0 0
3
4 3 1
3
1 4 2
3 1 3
1 2 1
1 0 0
0 0 1
0 0 0
2
1 3
1
2 1 3
0 1 1
1 1 0
0 0 1
1
3
3
1 2 1
1 1 3
1 3 3
1 1 1
0 1 1
1 1 0
2
5 4
3
1 2 3
1 3 3
1 4 1
1 0 1
1 1 1
0 1 1
4
4 4 2 5
2
3 2 2
4 1 2
1 1 1
0 1 0
0 1 0
3
1 2 4
0
1 1 0
1 0 0
1 1 0
3
3 4 5
3
3 5 3
1 1 1
2 3 3
1 1 0
1 1 1
0 0 0
3
5 4 1
1
1 4 2
1 1 1
1 0 1
1 0 1
1
4
2
1 3 1
1 4 3
1 1 0
1 0 0
0 1 0
1
1
1
1 1 2
1 0 0
0 0 0
0 1 1
2
2 5
1
2 2 1
0 0 0
1 0 1
0 0 1
2
1 4
0
1 1 1
1 1 0
0 1 0
3
5 5 5
1
2 1 3
0 1 1
0 0 1
1 1 0
2
1 3
0
0 0 0
1 0 1
1 1 1
1
5
0
1 0 0
1 0 0
0 1 0
2
5 4
2
1 3 2
2 4 2
0 1 1
1 1 0
0 0 1
1
1
0
0 0 0
0 0 1
1 1 0
4
4 3 5 4
2
4 4 1
2 1 3
0 0 1
1 1 0
1 0 1
2
3 2
0
1 0 1
1 1 0
0 1 1
2
1 3
2
2 3 1
1 1 1
0 0 0
1 1 1
0 1 1
3
1 4 2
1
1 1 1
1 1 0
1 0 1
1 0 0
1
5
2
1 1 1
1 4 1
0 1 1
1 1 0
0 0 0
4
2 5 2 1
2
1 1 3
2 5 1
0 0 1
1 1 1
1 0 0
3
4 2 5
0
0 1 1
0 1 1
1 0 0
1
4
1
1 4 3
0 1 0
1 0 0
0 1 1
1
1
1
1 1 3
0 1 0
0 1 0
0 0 1
3
1 4 5
1
2 1 2
0 0 1
1 1 1
0 0 0
1
1
0
0 1 1
1 0 0
1 0 1
2
3 3
0
0 0 0
1 0 0
1 1 1
1
3
2
1 1 2
1 3 2
1 0 0
0 1 0
0 0 0
2
4 3
3
2 1 1
2 3 3
2 2 1
0 1 1
1 0 0
0 1 1
1
1
0
1 0 1
0 1 1
1 0 0
4
1 1 5 5
1
4 1 2
0 1 1
1 1 1
0 0 1
3
5 1 5
1
2 1 1
0 0 1
1 1 0
1 0 0
1
1
1
1 1 2
0 0 0
1 1 1
0 1 1
2
5 2
2
2 2 1
2 1 1
0 1 1
1 0 0
0 1 0
2
1 5
3
2 2 2
2 5 3
1 1 2
1 0 1
0 1 0
0 0 1
2
4 2
2
1 3 3
1 1 2
0 0 1
0 1 1
0 1 0
1
2
2
1 1 1
1 2 1
0 1 0
0 0 0
1 1 0
4
1 2 4 5
3
4 1 1
3 2 3
1 1 2
1 0 1
1 0 0
1 1 0
3
4 2 3
3
3 3 1
3 1 2
2 2 1
1 0 0
1 0 1
0 0 1
1
5
3
1 2 3
1 4 2
1 1 1
0 1 0
0 1 1
1 0 0
2
2 4
0
1 0 1
1 0 0
0 1 1
1
2
0
1 1 1
1 0 0
1 0 1
2
4 1
2
1 4 2
1 3 1
1 1 1
0 0 0
1 0 0
4
1 4 2 2
0
1 1 1
0 1 0
0 1 0
1
2
1
1 2 1
1 1 0
0 1 0
1 0 0
4
3 3 3 4
3
2 3 1
3 1 3
1 3 2
0 0 1
0 0 0
0 1 0
2
3 1
1
1 3 2
1 0 1
1 1 0
0 1 1
4
4 3 2 5
0
1 0 0
1 0 1
1 1 0
4
3 1 4 1
1
4 1 1
0 0 1
0 1 1
1 1 1
2
1 3
2
1 1 2
2 3 2
0 0 0
1 0 0
1 0 1
1
1
1
1 1 3
1 0 0
0 0 0
0 1 1
1
3
0
1 1 1
1 1 1
0 1 1
4
4 2 2 3
2
2 2 1
1 4 2
1 0 0
0 1 0
1 1 0
2
4 3
3
1 2 3
1 4 3
2 2 2
0 0 1
1 1 1
1 0 0
1
3
0
0 1 0
0 1 0
1 1 0
4
3 1 3 1
2
3 2 2
1 1 2
0 0 0
0 1 0
1 1 0
2
4 1
2
1 3 2
1 4 2
0 0 0
0 0 0
1 1 0
3
3 5 4
3
1 1 3
1 2 3
3 2 1
1 0 0
1 0 1
1 1 1
3
2 3 4
3
2 3 3
1 1 3
1 2 3
0 1 0
0 0 0
0 1 1
4
2 4 1 3
2
3 1 3
1 1 1
1 1 1
1 1 0
1 0 1
3
4 3 2
2
3 1 3
1 1 2
0 1 1
1 1 1
1 1 0
2
5 4
0
1 1 1
0 0 0
1 0 0
2
2 5
1
2 1 2
1 0 1
0 0 0
0 1 1
3
2 1 3
3
3 3 1
2 1 1
1 1 1
0 1 1
1 1 1
0 0 1
2
1 4
1
1 1 1
0 1 0
0 0 1
0 0 1
1
3
1
1 2 1
0 1 0
0 0 1
1 0 1
1
3
3
1 3 3
1 1 1
1 2 3
1 0 0
1 1 0
0 0 1
1
5
3
1 2 3
1 5 1
1 1 3
0 1 1
0 1 0
1 0 0
3
4 3 5
3
2 3 3
2 1 2
1 4 2
0 0 1
0 0 1
0 0 1
4
2 3 3 1
2
1 1 3
4 1 1
1 1 0
1 1 1
0 0 1
1
5
2
1 4 2
1 1 1
1 1 0
1 0 1
0 1 1
1
4
2
1 4 3
1 2 2
1 1 1
1 0 1
1 0 1
4
1 2 1 4
2
1 1 2
4 1 1
0 1 0
1 0 0
0 0 0
3
3 3 4
1
1 2 2
1 0 1
1 1 1
0 0 1
4
3 5 1 1
3
3 1 2
1 2 2
4 1 1
0 1 1
1 0 0
1 0 1
3
3 2 5
1
2 2 1
1 1 0
0 0 1
0 1 0
1
5
2
1 3 2
1 2 1
0 0 1
0 1 1
0 0 0
3
3 2 3
3
1 1 3
1 3 2
3 2 2
1 1 0
1 1 0
1 0 0
4
1 3 4 2
0
0 1 0
1 1 0
1 0 1
2
1 1
0
1 1 0
1 0 0
1 1 1
3
3 2 5
2
2 1 3
1 3 3
0 0 0
1 0 0
0 0 1
2
3 2
1
1 1 3
1 0 1
1 1 0
0 1 0
3
5 4 5
0
1 0 0
0 1 1
1 1 1
2
2 4
1
1 1 2
0 1 0
0 1 0
1 1 1
1
3
2
1 3 1
1 2 2
1 1 0
1 1 1
0 1 1
1
4
0
0 1 1
0 1 1
0 1 0
1
3
1
1 2 2
0 1 1
1 1 0
0 1 0
2
1 1
1
2 1 3
0 0 1
0 0 0
0 1 1
2
1 1
1
2 1 3
0 0 1
1 1 0
1 0 1
`

type testCase struct {
	n     int
	a     []int
	fixed [][3]int
	fMat  [3][4]bool
}

// solve mirrors the logic from 1197F.go.
func solve(tc testCase) int {
	n := tc.n
	a := tc.a
	fixed := make([]map[int]int, n)
	for i := range fixed {
		fixed[i] = make(map[int]int)
	}
	for _, c := range tc.fixed {
		xi, yi, ci := c[0], c[1], c[2]
		fixed[xi-1][yi] = ci - 1
	}
	fMat := tc.fMat

	var transC [3][64]int
	for c := 0; c < 3; c++ {
		for idx := 0; idx < 64; idx++ {
			prev0 := idx >> 4
			prev1 := (idx >> 2) & 3
			prev2 := idx & 3
			used := [5]bool{}
			if fMat[c][1] {
				used[prev0] = true
			}
			if fMat[c][2] {
				used[prev1] = true
			}
			if fMat[c][3] {
				used[prev2] = true
			}
			mex := 0
			for used[mex] {
				mex++
			}
			newIdx := (mex << 4) | (prev0 << 2) | prev1
			transC[c][idx] = newIdx
		}
	}

	const S = 64
	mpow := make([][S][S]int, 31)
	for i := 0; i < S; i++ {
		for c := 0; c < 3; c++ {
			j := transC[c][i]
			mpow[0][i][j] = (mpow[0][i][j] + 1) % mod
		}
	}
	for k := 1; k < 31; k++ {
		for i := 0; i < S; i++ {
			for j := 0; j < S; j++ {
				sum := 0
				for t := 0; t < S; t++ {
					sum += mpow[k-1][i][t] * mpow[k-1][t][j]
					if sum >= mod*mod {
						sum -= mod * mod
					}
				}
				mpow[k][i][j] = sum % mod
			}
		}
	}

	dpXor := [4]int{1, 0, 0, 0}
	for i := 0; i < n; i++ {
		cnt := processStrip(a[i], fixed[i], fMat, transC, mpow)
		var newDp [4]int
		for x := 0; x < 4; x++ {
			sum := 0
			for g := 0; g < 4; g++ {
				sum = (sum + dpXor[x^g]*cnt[g]) % mod
			}
			newDp[x] = sum
		}
		dpXor = newDp
	}
	return dpXor[0]
}

// process one strip length L with fixed colors.
func processStrip(L int, fix map[int]int, fMat [3][4]bool, transC [3][64]int, mpow [][64][64]int) [4]int {
	const S = 64
	dpVec := make([]int, S)
	dpVec[0] = 1
	endInit := L
	if endInit > 3 {
		endInit = 3
	}
	for x := 1; x <= endInit; x++ {
		newVec := make([]int, S)
		for idx, val := range dpVec {
			if val == 0 {
				continue
			}
			prev0 := idx >> 4
			prev1 := (idx >> 2) & 3
			prev2 := idx & 3
			if c0, ok := fix[x]; ok {
				used := [5]bool{}
				if fMat[c0][1] && x-1 >= 1 {
					used[prev0] = true
				}
				if fMat[c0][2] && x-2 >= 1 {
					used[prev1] = true
				}
				if fMat[c0][3] && x-3 >= 1 {
					used[prev2] = true
				}
				mex := 0
				for used[mex] {
					mex++
				}
				ni := (mex << 4) | (prev0 << 2) | prev1
				newVec[ni] = (newVec[ni] + val) % mod
			} else {
				for c := 0; c < 3; c++ {
					used := [5]bool{}
					if fMat[c][1] && x-1 >= 1 {
						used[prev0] = true
					}
					if fMat[c][2] && x-2 >= 1 {
						used[prev1] = true
					}
					if fMat[c][3] && x-3 >= 1 {
						used[prev2] = true
					}
					mex := 0
					for used[mex] {
						mex++
					}
					ni := (mex << 4) | (prev0 << 2) | prev1
					newVec[ni] = (newVec[ni] + val) % mod
				}
			}
		}
		dpVec = newVec
	}
	if L <= 3 {
		var cnt [4]int
		for idx, val := range dpVec {
			if val == 0 {
				continue
			}
			g := idx >> 4
			cnt[g] = (cnt[g] + val) % mod
		}
		return cnt
	}
	cur := 3
	posList := make([]int, 0, len(fix))
	for pos := range fix {
		if pos >= 4 && pos <= L {
			posList = append(posList, pos)
		}
	}
	sort.Ints(posList)
	for _, pos := range posList {
		seg := pos - cur - 1
		if seg > 0 {
			dpVec = applyMatrix(dpVec, seg, mpow)
		}
		c0 := fix[pos]
		newVec := make([]int, S)
		for idx, val := range dpVec {
			if val == 0 {
				continue
			}
			ni := transC[c0][idx]
			newVec[ni] = (newVec[ni] + val) % mod
		}
		dpVec = newVec
		cur = pos
	}
	if L > cur {
		dpVec = applyMatrix(dpVec, L-cur, mpow)
	}
	var cnt [4]int
	for idx, val := range dpVec {
		if val == 0 {
			continue
		}
		g := idx >> 4
		cnt[g] = (cnt[g] + val) % mod
	}
	return cnt
}

// applyMatrix multiplies vector by M^seg.
func applyMatrix(vec []int, seg int, mpow [][64][64]int) []int {
	const S = 64
	res := make([]int, S)
	copy(res, vec)
	for k := 0; seg > 0; k++ {
		if seg&1 == 1 {
			tmp := make([]int, S)
			mat := mpow[k]
			for i := 0; i < S; i++ {
				if res[i] == 0 {
					continue
				}
				vi := res[i]
				for j := 0; j < S; j++ {
					tmp[j] = (tmp[j] + vi*mat[i][j]) % mod
				}
			}
			res = tmp
		}
		seg >>= 1
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	sc := bufio.NewScanner(strings.NewReader(testcasesRaw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("invalid test file")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, err
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("invalid test file")
		}
		n, _ := strconv.Atoi(sc.Text())
		a := make([]int, n)
		for j := 0; j < n; j++ {
			sc.Scan()
			a[j], _ = strconv.Atoi(sc.Text())
		}
		sc.Scan()
		cons, _ := strconv.Atoi(sc.Text())
		fixed := make([][3]int, cons)
		for j := 0; j < cons; j++ {
			sc.Scan()
			xi, _ := strconv.Atoi(sc.Text())
			sc.Scan()
			yi, _ := strconv.Atoi(sc.Text())
			sc.Scan()
			ci, _ := strconv.Atoi(sc.Text())
			fixed[j] = [3]int{xi, yi, ci}
		}
		var fMat [3][4]bool
		for r := 0; r < 3; r++ {
			for c := 1; c <= 3; c++ {
				sc.Scan()
				v, _ := strconv.Atoi(sc.Text())
				fMat[r][c] = (v == 1)
			}
		}
		tests = append(tests, testCase{n: n, a: a, fixed: fixed, fMat: fMat})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", len(tc.fixed))
	for _, f := range tc.fixed {
		fmt.Fprintf(&sb, "%d %d %d\n", f[0], f[1], f[2])
	}
	for r := 0; r < 3; r++ {
		fmt.Fprintf(&sb, "%d %d %d\n", btoi(tc.fMat[r][1]), btoi(tc.fMat[r][2]), btoi(tc.fMat[r][3]))
	}
	return sb.String()
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func runBinary(bin, input string) (string, error) {
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

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	expected := strconv.Itoa(solve(tc))
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
