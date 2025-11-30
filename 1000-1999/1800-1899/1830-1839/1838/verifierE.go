package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesE.txt.
const testcasesE = "MTAwCjIgMyAyIDIgMgoyIDEgMiAxIDIKMyAxIDMgMyAzIDIKMyAyIDIgMiAxIDEKMyAzIDMgMSAyIDIKMyAxIDYgMyA2IDQKMiAyIDMgMiAyCjEgMyAyIDIKMyAyIDUgNCA1IDMKNCA0IDMgMSAyIDIgMQoxIDEgNSAzCjQgMyAzIDMgMSAxIDIKMiA0IDQgMiAzCjQgMyA2IDIgMyAxIDEKMiAzIDYgNSAyCjEgMyAzIDIKNCAxIDIgMiAxIDIgMgoxIDMgNCAzCjIgNCA2IDYgMQozIDIgNSAzIDIgMwo0IDIgNCAxIDMgMSA0CjIgMyA0IDMgMQo0IDIgNSAyIDEgMSAxCjEgMiA2IDYKMiAxIDYgNCA1CjIgMyAyIDEgMgo0IDIgNSAyIDIgNCA0CjQgMSAzIDIgMiAxIDMKNCAyIDUgMiAxIDEgMwozIDIgNiAyIDIgNAozIDIgNCAxIDMgMQo0IDEgNSA0IDEgNCAyCjIgMyA0IDQgMwo0IDIgNCAzIDQgNCAxCjMgMiAyIDIgMSAyCjEgMiA1IDUKNCA0IDUgMiAxIDIgMgoxIDMgMiAyCjQgMiA2IDEgMiAyIDYKMyA0IDYgNCAxIDEKMSAxIDUgNQozIDIgMiAyIDEgMQozIDMgMiAxIDIgMgoyIDMgNSAyIDQKNCAxIDIgMSAyIDIgMgo0IDEgMyAzIDIgMiAyCjEgMiA0IDIKMiAyIDQgNCAzCjMgMSAzIDEgMyAyCjEgMSA1IDUKNCAxIDYgMiA0IDMgMwoyIDIgMiAxIDIKNCAzIDYgNCA1IDUgNAo0IDIgNCAzIDMgMiA0CjEgMiA2IDIKMyAzIDMgMyAyIDMKMSAxIDUgMgozIDMgNCAyIDMgMQoxIDEgNCAxCjIgMiA2IDQgNgoxIDEgMyAzCjQgMiA2IDYgMyA2IDQKMiA0IDMgMiAyCjIgMyA2IDYgNQo0IDQgNiA0IDMgMyA0CjQgMSA0IDIgNCAzIDQKNCAxIDMgMyAxIDIgMwozIDMgMiAyIDIgMQo0IDEgMyAxIDEgMiAxCjEgMiA1IDUKMyAxIDYgMiAyIDEKMiAxIDQgMSAxCjMgMSA0IDIgMiA0CjIgMSA2IDMgNgoxIDIgNiA0CjIgMiA0IDQgMQo0IDIgMyAyIDMgMiAzCjMgNCA0IDIgNCAxCjQgMSA0IDEgMyAyIDEKMiAxIDYgMSAyCjIgMSA1IDUgMgo0IDQgNCAyIDIgMSA0CjMgMyA1IDUgMSAyCjMgMyA1IDMgMSAyCjMgMSAyIDIgMiAyCjEgMiAyIDIKMiAxIDMgMyAxCjMgMiA0IDEgMSAxCjIgMyA2IDEgNgoxIDEgNSAyCjIgMyAzIDMgMwozIDMgNSA1IDQgMQoyIDQgNSAxIDUKMyAzIDMgMyAzIDIKMyAyIDMgMSAyIDIKMyAzIDUgNCAzIDUKMiAzIDQgMSAxCjQgMyA2IDMgNiA2IDIKMyAxIDMgMyAxIDMKMiAzIDMgMyAyCg=="

type testCase struct {
	n int
	m int
	k int
	arr []int
}

const mod int64 = 1_000_000_007
const maxN = 200000

var inv [maxN + 2]int64

func init() {
	inv[1] = 1
	for i := 2; i <= maxN+1; i++ {
		inv[i] = mod - (mod/int64(i))*inv[mod%int64(i)]%mod
	}
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

// Embedded solver logic from 1838E.go.
func solve(tc testCase) string {
	n := int64(tc.n)
	m := int64(tc.m)
	k := int64(tc.k)

	if k == 1 {
		return "1"
	}
	total := modPow(k, m)
	comb := int64(1)
	powTerm := modPow(k-1, m)
	sum := comb * powTerm % mod
	invK1 := modPow(k-1, mod-2)
	limit := n - 1
	if limit > m {
		limit = m
	}
	for i := int64(1); i <= limit; i++ {
		comb = comb * ((m - i + 1) % mod) % mod
		comb = comb * inv[i] % mod
		powTerm = powTerm * invK1 % mod
		sum = (sum + comb*powTerm%mod) % mod
	}
	ans := (total - sum) % mod
	if ans < 0 {
		ans += mod
	}
	return strconv.FormatInt(ans, 10)
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesE)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(sc.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d n: %v", i+1, err)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing m", i+1)
		}
		m, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d m: %v", i+1, err)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("case %d missing k", i+1)
		}
		k, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d k: %v", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d missing a[%d]", i+1, j)
			}
			arr[j], err = strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d a[%d]: %v", i+1, j, err)
			}
		}
		cases = append(cases, testCase{n: n, m: m, k: k, arr: arr})
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d %d %d\n", tc.n, tc.m, tc.k)
		for i, v := range tc.arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		want := solve(tc)
		got, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", idx+1, input.String(), want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
