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

type testCase struct {
	n int
	a []int
}

func validSegments(a []int, segs [][2]int) bool {
	if len(segs) == 0 {
		return false
	}
	n := len(a)
	if segs[0][0] != 1 {
		return false
	}
	if segs[len(segs)-1][1] != n {
		return false
	}
	for i := 0; i < len(segs); i++ {
		l := segs[i][0]
		r := segs[i][1]
		if l > r || l < 1 || r > n {
			return false
		}
		if i > 0 && segs[i-1][1]+1 != l {
			return false
		}
		sum := 0
		for j := l - 1; j < r; j++ {
			sum += a[j]
		}
		if sum == 0 {
			return false
		}
	}
	return true
}

func hasSolution(a []int) bool {
	allZero := true
	total := 0
	for _, v := range a {
		if v != 0 {
			allZero = false
		}
		total += v
	}
	if allZero {
		return false
	}
	if total != 0 {
		return true
	}
	prefix := 0
	for i := 0; i < len(a)-1; i++ {
		prefix += a[i]
		if prefix != 0 {
			return true
		}
	}
	return false
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintln(&sb, tc.n)
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	token := scanner.Text()
	if token == "NO" {
		if hasSolution(tc.a) {
			return fmt.Errorf("reported NO but solution exists")
		}
		if scanner.Scan() {
			return fmt.Errorf("extra output")
		}
		return nil
	}
	if token != "YES" {
		return fmt.Errorf("first token should be YES or NO")
	}
	if !scanner.Scan() {
		return fmt.Errorf("missing k")
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil || k <= 0 {
		return fmt.Errorf("invalid k")
	}
	segs := make([][2]int, k)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing l for segment %d", i+1)
		}
		l, err1 := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			return fmt.Errorf("missing r for segment %d", i+1)
		}
		r, err2 := strconv.Atoi(scanner.Text())
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid segment %d", i+1)
		}
		segs[i] = [2]int{l, r}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	if !hasSolution(tc.a) {
		return fmt.Errorf("reported YES but no valid solution exists")
	}
	if !validSegments(tc.a, segs) {
		return fmt.Errorf("invalid segments")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{1, []int{0}},
		{1, []int{5}},
		{3, []int{1, -1, 0}},
		{3, []int{1, 2, 3}},
		{4, []int{0, 0, 5, -5}},
		{2, []int{-1, 1}},
		{2, []int{1, -1}},
		{5, []int{0, 0, 1, 0, 0}},
		{3, []int{-5, 5, 0}},
		{3, []int{0, 0, 1}},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(21) - 10
		}
		cases = append(cases, testCase{n, a})
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\n", i+1, err, tc.a)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
