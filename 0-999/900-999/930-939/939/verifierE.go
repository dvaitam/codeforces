package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type test struct {
	in  string
	out string
}

type op struct {
	typ int
	val int64
}

func solveE(ops []op) string {
	// Recompute the expected answer per query from scratch.
	// Maintain prefix sums and for each type-2 query compute
	// min over k in [0..idx-1] of (P[k] + a[idx])/(k+1).
	Q := len(ops)
	a := make([]int64, Q+5)    // 1-indexed values
	pref := make([]int64, Q+5) // prefix sums, pref[0]=0
	idx := 0
	results := make([]string, 0)
	for _, o := range ops {
		if o.typ == 1 {
			idx++
			a[idx] = o.val
			pref[idx] = pref[idx-1] + o.val
		} else {
			if idx == 0 {
				results = append(results, fmt.Sprintf("%.8f", 0.0))
				continue
			}
			last := float64(a[idx])
			minAvg := last // k=0 gives avg=last
			for k := 1; k <= idx-1; k++ {
				avg := float64(pref[k]+a[idx]) / float64(k+1)
				if avg < minAvg {
					minAvg = avg
				}
			}
			diff := last - minAvg
			results = append(results, fmt.Sprintf("%.8f", diff))
		}
	}
	return strings.Join(results, "\n")
}

func generateTests() []test {
	rand.Seed(5)
	tests := make([]test, 0, 100)
	for i := 0; i < 100; i++ {
		Q := rand.Intn(10) + 1
		ops := make([]op, 0, Q)
		hasQuery := false
		inserted := 0
		for j := 0; j < Q; j++ {
			if inserted == 0 || rand.Intn(2) == 0 {
				v := rand.Int63n(100) + 1
				ops = append(ops, op{typ: 1, val: v})
				inserted++
			} else {
				ops = append(ops, op{typ: 2})
				hasQuery = true
			}
		}
		if !hasQuery {
			ops = append(ops, op{typ: 2})
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(ops))
		for _, o := range ops {
			if o.typ == 1 {
				fmt.Fprintf(&sb, "1 %d\n", o.val)
			} else {
				sb.WriteString("2\n")
			}
		}
		tests = append(tests, test{in: sb.String(), out: solveE(ops)})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := run(bin, t.in)
		if err != nil {
			fmt.Printf("Test %d failed to run: %v\n", i+1, err)
			fmt.Print(out)
			os.Exit(1)
		}

		// Tolerant float comparison (accepts different formatting)
		expLines := strings.Fields(strings.TrimSpace(t.out))
		gotFields := strings.Fields(strings.TrimSpace(out))
		if len(expLines) != len(gotFields) {
			fmt.Printf("Test %d failed. Expected %d lines, got %d. Input:\n%sExpected:\n%s\nGot:\n%s\n", i+1, len(expLines), len(gotFields), t.in, t.out, out)
			os.Exit(1)
		}
		ok := true
		const eps = 1e-6
		for j := range expLines {
			expV, err1 := strconv.ParseFloat(expLines[j], 64)
			gotV, err2 := strconv.ParseFloat(gotFields[j], 64)
			if err1 != nil || err2 != nil || (gotV-expV > eps) || (expV-gotV > eps) {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Printf("Test %d failed. Expected %q, got %q. Input:\n%s", i+1, t.out, out, t.in)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
