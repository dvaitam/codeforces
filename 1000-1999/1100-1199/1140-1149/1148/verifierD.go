package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin string, in []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	ref := "./refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1148D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genCase() (string, [][2]int) {
	n := rand.Intn(5) + 2 // 2..6
	vals := rand.Perm(2 * n)
	pairs := make([][2]int, n)
	for i := 0; i < n; i++ {
		a := vals[2*i] + 1
		b := vals[2*i+1] + 1
		if rand.Intn(2) == 0 {
			pairs[i] = [2]int{a, b}
		} else {
			pairs[i] = [2]int{b, a}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", pairs[i][0], pairs[i][1]))
	}
	return sb.String(), pairs
}

func parseIndices(out string, n int) ([]int, error) {
	fields := strings.Fields(strings.TrimSpace(out))
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty output")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid count")
	}
	if len(fields)-1 != t {
		return nil, fmt.Errorf("expected %d indices", t)
	}
	idx := make([]int, t)
	for i := 0; i < t; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil || v < 1 || v > n {
			return nil, fmt.Errorf("bad index")
		}
		idx[i] = v
	}
	return idx, nil
}

func isGood(seq []int) bool {
	if len(seq) < 2 {
		return true
	}
	// try pattern < > < > ...
	ok := true
	less := seq[0] < seq[1]
	if seq[0] == seq[1] {
		ok = false
	}
	expectLess := !less
	for i := 1; i < len(seq)-1 && ok; i++ {
		if expectLess {
			if !(seq[i] < seq[i+1]) {
				ok = false
			}
		} else {
			if !(seq[i] > seq[i+1]) {
				ok = false
			}
		}
		expectLess = !expectLess
	}
	if ok {
		return true
	}
	// try opposite pattern
	ok = true
	less = seq[0] > seq[1]
	if seq[0] == seq[1] {
		return false
	}
	expectLess = !less
	for i := 1; i < len(seq)-1 && ok; i++ {
		if expectLess {
			if !(seq[i] < seq[i+1]) {
				ok = false
			}
		} else {
			if !(seq[i] > seq[i+1]) {
				ok = false
			}
		}
		expectLess = !expectLess
	}
	return ok
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		input, pairs := genCase()
		wantOut, err := run(ref, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		refIdx, err := parseIndices(wantOut, len(pairs))
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad reference output on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotOut, err := run(bin, []byte(input))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		candIdx, err := parseIndices(gotOut, len(pairs))
		if err != nil {
			fmt.Fprintf(os.Stderr, "wrong format on test %d: %v\ninput:\n%soutput:\n%s", t+1, err, input, gotOut)
			os.Exit(1)
		}
		if len(candIdx) != len(refIdx) {
			fmt.Fprintf(os.Stderr, "wrong subset size on test %d\n", t+1)
			os.Exit(1)
		}
		seq := make([]int, 0, 2*len(candIdx))
		for _, id := range candIdx {
			pr := pairs[id-1]
			seq = append(seq, pr[0], pr[1])
		}
		if !isGood(seq) {
			fmt.Fprintf(os.Stderr, "sequence not good on test %d\ninput:\n%soutput:\n%s", t+1, input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
