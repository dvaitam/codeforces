package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE1.txt so the verifier is self-contained.
const testcasesRaw = `10 0 33 95 46 89 95 84 68 4 60 100
4 0 84 7 21 15
6 0 61 32 49 70 14 74
4 0 2 94 28 53
5 0 24 99 50 21 98
2 0 18 80
10 0 57 17 17 1 1 27 100 28 22 22
5 0 41 26 70 87 81
4 0 24 89 26 50
5 0 3 47 54 22 19
5 0 9 43 39 78 76
1 0 77
6 0 9 40 46 40 62 90
6 0 24 62 61 91 23 8
5 0 3 96 46 52 3
9 0 54 47 49 75 2 58 6 91 24
10 0 26 16 97 32 60 45 66 46 68 33
8 0 14 76 96 100 48 38 5 56
2 0 27 44
9 0 79 47 19 44 36 90 70 12 40
6 0 40 23 11 81 20 93
5 0 62 21 93 7 11
10 0 69 52 5 31 95 77 45 33 59 84
7 0 19 8 82 5 64 43 27
3 0 94 73 17
7 0 14 22 56 48 20 8 54
5 0 19 59 80 22 67
8 0 63 89 94 41 62 36 38 61
7 0 19 15 49 69 23 81 64
6 0 24 12 63 35 66 71
9 0 47 9 100 46 89 76 85 5 98
5 0 47 72 91 86 36
8 0 34 99 89 92 38 44 84 23
10 0 2 61 71 100 33 42 86 36 60 37
9 0 83 87 46 45 36 83 45 95 53
6 0 23 89 58 47 43 67
3 0 68 22 26
6 0 62 37 89 11 93 86
7 0 22 79 100 75 67 86 54
5 0 80 71 100 82 35
1 0 26
3 0 76 57 80
3 0 29 98 88
3 0 81 92 6
8 0 29 22 7 18 15 41 24 62
4 0 71 5 54 60
6 0 49 85 79 10 76 27
4 0 92 48 1 45
7 0 36 53 15 89 71 48 5
9 0 79 39 13 38 70 66 44 75 38
6 0 17 54 53 73 83 69
6 0 60 19 21 77 49 73
8 0 26 18 78 12 45 85 1 49
2 0 42 73
10 0 70 19 42 81 73 49 55 56 29 64
5 0 62 91 49 50 21
10 0 77 34 95 39 64 33 54 3 41 40
8 0 37 19 62 4 16 85 80 57
4 0 38 6 18 51
1 0 62
9 0 72 36 32 61 5 32 63 35 20
5 0 38 64 78 61 67
10 0 96 16 3 98 17 39 37 69 91 44
10 0 38 94 68 4 60 45 47 88 96 76
3 0 5 1 33
9 0 59 88 14 88 70 25 2 55 100
7 0 77 74 89 91 81 84 62
7 0 61 51 88 93 26 38 60
2 0 39 1
7 0 75 37 83 100 61 40 19
3 0 62 89 71
8 0 43 69 20 55 75 70 7 9
4 0 35 11 9 85
1 0 43
7 0 9 52 90 63 7 16 16
4 0 79 83 15 92
3 0 38 91 57
3 0 24 79 24
7 0 21 9 80 28 6 72 14
7 0 96 10 36 8 74 74 16
7 0 80 18 2 56 12 41 88
10 0 63 63 46 84 48 8 18 90 38 20
10 0 81 87 65 38 72 71 80 29 34 9
9 0 31 33 97 37 67 18 31 48 59
7 0 23 17 92 3 84 44 11
10 0 86 5 12 16 65 77 59 31 50 60
8 0 42 14 68 4 70 93 50 7
3 0 55 88 29
2 0 11 86
8 0 28 18 90 80 49 46 31 38
6 0 79 91 45 50 49 18
6 0 83 38 81 56 47 67
1 0 76
10 0 28 95 24 51 9 13 5 5 24 26
4 0 6 63 62 85
6 0 1 55 61 39 80 55
6 0 60 60 13 25 20 84
3 0 10 48 50
8 0 20 71 33 15 36 21 97 37
4 0 5 62 5 45
`

type testCase struct {
	n   int
	k   int
	arr []int
}

func buildPrimes(limit int) []int {
	isComp := make([]bool, limit+1)
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if !isComp[i] {
			primes = append(primes, i)
			if i*i <= limit {
				for j := i * i; j <= limit; j += i {
					isComp[j] = true
				}
			}
		}
	}
	return primes
}

func canonical(x int, primes []int) int {
	res := 1
	for _, p := range primes {
		if p*p > x {
			break
		}
		cnt := 0
		for x%p == 0 {
			x /= p
			cnt ^= 1
		}
		if cnt == 1 {
			res *= p
		}
	}
	if x > 1 {
		res *= x
	}
	return res
}

func segmentsFor(arr []int, mask int) int {
	seen := map[int]bool{}
	seg := 1
	for i, v := range arr {
		if mask>>i&1 == 1 {
			v = -100000 - i
		}
		if seen[v] {
			seg++
			seen = map[int]bool{}
		}
		seen[v] = true
	}
	return seg
}

func solveBrute(arr []int, k int) int {
	n := len(arr)
	best := n
	for mask := 0; mask < 1<<n; mask++ {
		if bits.OnesCount(uint(mask)) > k {
			continue
		}
		seg := segmentsFor(arr, mask)
		if seg < best {
			best = seg
		}
	}
	return best
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var res []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		k, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %v", idx+1, err)
		}
		if len(parts) != n+2 {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n+2, len(parts))
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse value %d: %v", idx+1, i+1, err)
			}
			arr[i] = v
		}
		res = append(res, testCase{n: n, k: k, arr: arr})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	primes := buildPrimes(3200)

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for i, tc := range cases {
		arr := make([]int, tc.n)
		for idx, v := range tc.arr {
			arr[idx] = canonical(v, primes)
		}
		expected := strconv.Itoa(solveBrute(arr, tc.k))

		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d %d\n", tc.n, tc.k)
		for idx, v := range tc.arr {
			if idx > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Printf("test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
