package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesC.txt.
const testcasesC = "MiAwCjMgMiAxIDMgMyAzCjYgMSA1IDUKNCAzIDQgNCA0IDQgMSAxCjcgMyAzIDYgNCA1IDUgNQo1IDEgMSAyCjcgMSAyIDYKMTAgMiA5IDkgOCA5CjEwIDIgMTAgMTAgNiA5CjQgMyA0IDQgNCA0IDQgNAo5IDMgNiA5IDggOCA2IDcKNiAzIDMgNSA2IDYgMyA0CjkgMiAyIDcgMSA0CjMgMAoyIDEgMSAxCjEwIDEgNSA2CjUgMAo4IDAKMiAxIDIgMgo1IDAKMyAwCjMgMAoyIDIgMSAyIDIgMgo0IDEgMSA0CjIgMAo0IDAKMiAxIDEgMgo3IDMgMSAzIDQgNCAzIDYKNCAzIDIgMiAzIDMgMSA0CjQgMyA0IDQgMiAzIDMgNAo4IDAKMTAgMSAxIDUKMiAwCjQgMSAxIDQKNSAwCjUgMSA0IDQKNiAwCjUgMiAzIDUgNCA1CjEwIDAKNCAwCjggMyAzIDMgMiAzIDIgMgoyIDAKNSAwCjUgMAoxMCAzIDggOSA5IDEwIDQgOQo1IDMgNCA0IDUgNSA0IDQKMyAzIDIgMiAzIDMgMyAzCjYgMiAzIDMgNiA2CjMgMAo2IDEgNiA2CjkgMAo4IDMgOCA4IDIgMiA1IDUKNyAyIDYgNiAyIDUKNSAwCjcgMyA2IDcgMiA0IDQgNAo2IDAKMyAwCjcgMyAyIDcgMSAxIDUgNwo5IDAKOSAyIDYgOSAzIDkKNyAyIDQgNyA2IDcKOSAyIDcgNyAzIDYKNiAzIDYgNiA1IDUgMSAzCjQgMSA0IDQKMyAwCjQgMiA0IDQgMyA0CjQgMiAxIDIgNCA0CjcgMSA2IDcKNCAxIDQgNAo4IDIgMyA2IDggOAo4IDEgNyA3CjkgMiA3IDggNyA5CjkgMiA5IDkgMiA1CjEwIDEgNyAxMAoyIDEgMiAyCjQgMAoyIDEgMSAyCjUgMAo4IDEgNSA2CjkgMSAxIDcKOSAyIDkgOSA4IDgKMyAyIDEgMiAzIDMKOSAyIDggOSA5IDkKMiAwCjcgMSA3IDcKNiAxIDEgMgo5IDMgOCA5IDMgMyA1IDkKOSAwCjcgMAoxMCAzIDEwIDEwIDQgMTAgNSA4CjQgMyAzIDMgMyA0IDMgNAo2IDMgNSA1IDUgNSA0IDYKMTAgMSA5IDkKNiAwCjUgMyA1IDUgNSA1IDEgMQozIDMgMiAyIDIgMyAxIDIKOSAyIDMgNiA4IDkKOSAxIDggOAo2IDIgMiAyIDIgNQo1IDIgMiA0IDIgMwo="

type interval struct{ l, r int }

type testCase struct {
	n, k      int
	intervals []interval
}

const mod int64 = 998244353
const maxN = 300000

var fact [maxN + 1]int64
var inv [maxN + 1]int64

func powmod(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func init() {
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	inv[maxN] = powmod(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		inv[i-1] = inv[i] * int64(i) % mod
	}
}

func catalan(n int) int64 {
	if n < 0 {
		return 0
	}
	return fact[2*n] * inv[n] % mod * inv[n+1] % mod
}

// Embedded solver logic from 1830C.go.
func solve(tc testCase) string {
	n, k := tc.n, tc.k
	nodes := make([]interval, 0, k+1)
	ok := true
	for _, iv := range tc.intervals {
		if (iv.r-iv.l+1)%2 == 1 {
			ok = false
		}
		nodes = append(nodes, iv)
	}
	if n%2 == 1 {
		ok = false
	}
	if !ok {
		return "0"
	}
	nodes = append(nodes, interval{1, n})
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].l == nodes[j].l {
			return nodes[i].r > nodes[j].r
		}
		return nodes[i].l < nodes[j].l
	})
	// deduplicate consecutive duplicates
	uniq := nodes[:0]
	for _, iv := range nodes {
		if len(uniq) == 0 || iv != uniq[len(uniq)-1] {
			uniq = append(uniq, iv)
		}
	}
	nodes = uniq

	children := make([][]int, len(nodes))
	stack := []int{}
	valid := true
	for idx, iv := range nodes {
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if nodes[top].l <= iv.l && iv.r <= nodes[top].r {
				break
			}
			stack = stack[:len(stack)-1]
		}
		if len(stack) > 0 && iv.r > nodes[stack[len(stack)-1]].r {
			valid = false
			break
		}
		if len(stack) > 0 {
			p := stack[len(stack)-1]
			children[p] = append(children[p], idx)
		}
		stack = append(stack, idx)
	}
	if !valid {
		return "0"
	}

	for v := range children {
		sort.Slice(children[v], func(i, j int) bool {
			return nodes[children[v][i]].l < nodes[children[v][j]].l
		})
	}

	rootIdx := -1
	for i, iv := range nodes {
		if iv.l == 1 && iv.r == n {
			rootIdx = i
			break
		}
	}
	if rootIdx == -1 {
		return "0"
	}

	ansVal := make([]int64, len(nodes))
	post := []int{}
	st2 := []int{rootIdx}
	seen := make([]bool, len(nodes))
	for len(st2) > 0 {
		v := st2[len(st2)-1]
		st2 = st2[:len(st2)-1]
		if seen[v] {
			post = append(post, v)
			continue
		}
		seen[v] = true
		st2 = append(st2, v)
		for _, c := range children[v] {
			st2 = append(st2, c)
		}
	}

	for i := 0; i < len(post); i++ {
		v := post[i]
		used := 0
		cur := int64(1)
		for _, c := range children[v] {
			cur = cur * ansVal[c] % mod
			used += (nodes[c].r - nodes[c].l + 1) / 2
		}
		free := (nodes[v].r-nodes[v].l+1)/2 - used
		if free < 0 {
			cur = 0
		} else {
			cur = cur * catalan(free) % mod
		}
		ansVal[v] = cur
	}
	return strconv.FormatInt(ansVal[rootIdx], 10)
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesC)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(bytes.NewReader(raw))
	sc.Split(bufio.ScanWords)
	cases := []testCase{}
	for sc.Scan() {
		n, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %v", err)
		}
		if !sc.Scan() {
			return nil, fmt.Errorf("missing k")
		}
		k, err := strconv.Atoi(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("parse k: %v", err)
		}
		tmp := make([]interval, k)
		for i := 0; i < k; i++ {
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d l missing", len(cases)+1)
			}
			l, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d l: %v", len(cases)+1, err)
			}
			if !sc.Scan() {
				return nil, fmt.Errorf("case %d r missing", len(cases)+1)
			}
			r, err := strconv.Atoi(sc.Text())
			if err != nil {
				return nil, fmt.Errorf("case %d r: %v", len(cases)+1, err)
			}
			tmp[i] = interval{l: l, r: r}
		}
		seen := make(map[interval]bool)
		intervals := make([]interval, 0, len(tmp))
		for _, iv := range tmp {
			if iv.l == 1 && iv.r == n {
				continue
			}
			if !seen[iv] {
				seen[iv] = true
				intervals = append(intervals, iv)
			}
		}
		cases = append(cases, testCase{n: n, k: len(intervals), intervals: intervals})
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
		fmt.Fprintf(&input, "1\n%d %d\n", tc.n, tc.k)
		for _, iv := range tc.intervals {
			fmt.Fprintf(&input, "%d %d\n", iv.l, iv.r)
		}
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
