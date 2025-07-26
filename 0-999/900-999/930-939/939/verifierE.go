package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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
	Q := len(ops)
	a := make([]int64, Q+5)
	var idx, cnt int
	var sum int64
	var Max int64
	var aver float64 = 1e18
	results := make([]string, 0)
	for _, o := range ops {
		if o.typ == 1 {
			idx++
			a[idx] = o.val
			Max = o.val
		} else {
			aver = float64(sum+a[idx]) / float64(cnt+1)
			for cnt < idx-1 {
				tmp := float64(sum+a[idx]+a[cnt+1]) / float64(cnt+2)
				if tmp < aver {
					aver = tmp
					sum += a[cnt+1]
					cnt++
				} else {
					break
				}
			}
			diff := float64(Max) - aver
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
		if strings.TrimSpace(out) != t.out {
			fmt.Printf("Test %d failed. Expected %q, got %q. Input:\n%s", i+1, t.out, out, t.in)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
