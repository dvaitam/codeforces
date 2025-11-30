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

const testcasesBase64 = "MTAwCjQgMCA0IDMgMiAxCjIgMiAyIDEgMiAxCjEwIDIgNiAxIDkgNCA4IDcgMyAyIDUgMTAgNyAxCjIgMiAyIDEgMSAyCjEwIDUgNiAxMCA1IDkgMiA3IDQgOCAxIDMgOCA0IDMgOSAxMAo3IDMgNSA2IDEgMyAyIDQgNyAzIDEgNAo1IDMgNSA0IDEgMiAzIDQgMSA1CjggMiAxIDggMyA0IDYgNyA1IDIgMyAxCjYgMSA1IDIgNCAxIDYgMyA0CjYgMiA1IDQgNiAyIDEgMyAzIDQKMyAxIDMgMSAyIDMKOSA0IDUgMyA5IDcgOCA2IDQgMSAyIDkgMiA1IDQKMyAwIDMgMSAyCjUgMCAzIDIgNCAxIDUKNyAyIDIgMSA2IDcgMyA0IDUgNCA3CjUgMSAzIDUgNCAxIDIgMQoyIDEgMSAyIDIKOCAwIDUgNyAxIDYgNCAyIDggMwozIDMgMyAxIDIgMSAzIDIKNiAxIDEgNSA0IDMgNiAyIDQKMTAgMTAgNiAxMCA0IDIgMSAzIDcgNSA5IDggOSA4IDEwIDQgMiA1IDYgMyAxIDcKNyA1IDQgMSA3IDYgMiAzIDUgMiAzIDQgNiA1CjEgMCAxCjkgMCA0IDUgNyA4IDYgOSAxIDIgMwo0IDEgMSA0IDIgMyAzCjQgMyAyIDQgMyAxIDIgNCAxCjcgMiA1IDMgNCA2IDEgMiA3IDQgNgo2IDQgMSA1IDIgMyA2IDQgMyA2IDIgNQo2IDMgNSA0IDYgMiAxIDMgNiAyIDEKMyAxIDMgMSAyIDMKMTAgNSA2IDEwIDUgOSAyIDcgNCA4IDEgMyA4IDQgMyA5IDEwCjc gMyA1IDYgMSAzIDIgNCA3IDMgMSA0CjUgMyA1IDQgMSAyIDMgNCAxIDUKOCAyIDEgOCAzIDQgNiA3IDUgMiAzIDEKNiAxIDUgMiA0IDEgNiAzIDQKNiAyIDUgNCA2IDIgMSAzIDMgNAozIDEgMyAxIDIgMwo5IDQgNSA0IDggNiA0IDEgMiAyIDkgMiA1IDQKMyAwIDMgMSAyCjUgMCAzIDIgNCAxIDUKNyAyIDIgMSA2IDcgMyA0IDUgNCA3CjUgMSAzIDUgNCAxIDIgMQoyIDEgMSAyIDIKOCAwIDUgNyAxIDYgNCAyIDggMwozIDMgMyAxIDIgMSAzIDIKNiAxIDEgNSA0IDMgNiAyIDQKMTAgMTAgNiAxMCA0IDIgMSAzIDcgNSA5IDggOSA4IDEwIDQgMiA1IDYgMyAxIDcKNyA1IDQgMSA3IDYgMiAzIDUgMiAzIDQgNiA1CjEgMCAxCjkgMCA0IDUgNyA4IDYgOSAxIDIgMwo0IDEgMSA0IDIgMyAzCjQgMyAyIDQgMyAxIDIgNCAxCjcgMiA1IDMgNCA2IDEgMiA3IDQgNgo2IDQgMSA1IDIgMyA2IDQgMyA2IDIgNQo2IDMgNSA0IDYgMiAxIDMgNiAyIDEKMyAxIDMgMSAyIDMKOSA4IDUgMiAzIDUgNiA3IDUgMiA2IDQgMSAyIDkgMiA1IDgKNyA1IDMgNCA3IDggMiAzIDUgNyAyCjQgMyAxIDYgMyAyCjUgMCAxIDMgNyAyIDYgNSA0CjEgMCAxCjkgMCA0IDUgNyA4IDYgOSAxIDIgMwo0IDEgMSA0IDIgMyAzCjQgMyAyIDQgMyAxIDIgNCAxCjcgMiA1IDMgNCA2IDEgMiA3IDQgNgo2IDQgMSA1IDIgMyA2IDQgMyA2IDIgNQo2IDMgNSA0IDYgMiAxIDMgNiAyIDEKMyAxIDMgMSAyIDMKMiAyIDEgMiAxIDI KNjkK"

type testCase struct {
	n int
	k int
	a []int
	b []int
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesData))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %w", err)
	}
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			break
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			break
		}
		k, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return cases, nil
			}
			a[j], _ = strconv.Atoi(scan.Text())
		}
		b := make([]int, k)
		for j := 0; j < k; j++ {
			if !scan.Scan() {
				return cases, nil
			}
			b[j], _ = strconv.Atoi(scan.Text())
		}
		cases = append(cases, testCase{n: n, k: k, a: a, b: b})
	}
	return cases, nil
}

type bit struct {
	n    int
	data []int
}

func newBIT(n int) *bit {
	return &bit{n: n, data: make([]int, n+1)}
}

func (b *bit) add(i, v int) {
	for i <= b.n {
		b.data[i] += v
		i += i & -i
	}
}

func (b *bit) sum(i int) int {
	res := 0
	for i > 0 {
		res += b.data[i]
		i -= i & -i
	}
	return res
}

func (b *bit) findByOrder(k int) int {
	idx := 0
	bitMask := 1
	for bitMask<<1 <= b.n {
		bitMask <<= 1
	}
	for d := bitMask; d > 0; d >>= 1 {
		next := idx + d
		if next <= b.n && b.data[next] < k {
			idx = next
			k -= b.data[next]
		}
	}
	return idx + 1
}

type kv struct {
	val int
	pos int
}

func referenceSolve(tc testCase) (string, error) {
	pos := make([]int, tc.n+1)
	for i, v := range tc.a {
		pos[v] = i + 1
	}
	isKeeper := make([]bool, tc.n+1)
	bElems := make([]kv, 0, tc.k)
	for _, v := range tc.b {
		isKeeper[v] = true
		bElems = append(bElems, kv{val: v, pos: pos[v]})
	}
	remElems := make([]kv, 0, tc.n-tc.k)
	for v := 1; v <= tc.n; v++ {
		if !isKeeper[v] {
			remElems = append(remElems, kv{val: v, pos: pos[v]})
		}
	}
	sort.Slice(bElems, func(i, j int) bool { return bElems[i].val < bElems[j].val })

	remBIT := newBIT(tc.n)
	for i := 1; i <= tc.n; i++ {
		remBIT.add(i, 1)
	}
	blockerBIT := newBIT(tc.n)

	var ans int64
	idxB := 0
	for _, item := range remElems {
		x := item.val
		px := item.pos
		for idxB < len(bElems) && bElems[idxB].val < x {
			blockerBIT.add(bElems[idxB].pos, 1)
			idxB++
		}
		sumLeft := blockerBIT.sum(px - 1)
		leftBound := 0
		if sumLeft > 0 {
			leftBound = blockerBIT.findByOrder(sumLeft)
		}
		sumBefore := blockerBIT.sum(px)
		totalBl := blockerBIT.sum(tc.n)
		rightBound := tc.n + 1
		if sumBefore < totalBl {
			rightBound = blockerBIT.findByOrder(sumBefore + 1)
		}
		l := leftBound + 1
		r := rightBound - 1
		if l <= r {
			count := remBIT.sum(r) - remBIT.sum(leftBound)
			ans += int64(count)
		}
		remBIT.add(px, -1)
	}
	return strconv.FormatInt(ans, 10), nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierE /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected, err := referenceSolve(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for i, v := range tc.b {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		out, stderr, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
