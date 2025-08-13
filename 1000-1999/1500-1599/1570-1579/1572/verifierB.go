package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseB struct {
	n   int
	arr []int
}

func generateCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(8) + 3 // 3..10
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(2)
	}
	return testCaseB{n: n, arr: arr}
}

func buildInputB(t testCaseB) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(t.n))
	sb.WriteString("\n")
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteString("\n")
	return sb.String()
}

func solveB(reader *bufio.Reader) string {
	var T int
	fmt.Fscan(reader, &T)
	out := strings.Builder{}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, n+1)
		pos := 0
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			a[i] ^= a[i-1]
			if i%2 == 1 && a[i] == 0 {
				pos = i
			}
		}
		if pos != 0 && a[n] == 0 {
			out.WriteString("YES\n")
			m := n - 1
			if n != pos {
				m--
			}
			out.WriteString(fmt.Sprintln(m))
			for i := pos + 1; i+2 <= n; i += 2 {
				out.WriteString(fmt.Sprintf("%d ", i))
			}
			for i := pos; ; {
				i -= 2
				if i < 1 {
					break
				}
				out.WriteString(fmt.Sprintf("%d ", i))
			}
			for i := 1; i+2 <= n; i += 2 {
				out.WriteString(fmt.Sprintf("%d ", i))
			}
			out.WriteByte('\n')
		} else {
			out.WriteString("NO\n")
		}
	}
	return strings.TrimSpace(out.String())
}

func expectedB(t testCaseB) string {
	input := buildInputB(t)
	return solveB(bufio.NewReader(strings.NewReader(input)))
}

func isSolvableB(t testCaseB) bool {
	exp := strings.TrimSpace(expectedB(t))
	return strings.HasPrefix(exp, "YES")
}

func verifyOutputB(t testCaseB, output string) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("empty output")
	}
	tokens := strings.Fields(output)
	solvable := isSolvableB(t)
	if tokens[0] == "NO" {
		if solvable {
			return fmt.Errorf("expected YES, got NO")
		}
		if len(tokens) != 1 {
			return fmt.Errorf("unexpected tokens after NO")
		}
		return nil
	}
	if tokens[0] != "YES" {
		return fmt.Errorf("first token should be YES or NO")
	}
	if !solvable {
		return fmt.Errorf("expected NO, got YES")
	}
	if len(tokens) < 2 {
		return fmt.Errorf("missing k")
	}
	k, err := strconv.Atoi(tokens[1])
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	if k < 0 || k > t.n {
		return fmt.Errorf("k out of range")
	}
	if len(tokens) != 2+k {
		return fmt.Errorf("expected %d operations, got %d", k, len(tokens)-2)
	}
	arr := append([]int(nil), t.arr...)
	for _, opStr := range tokens[2:] {
		op, err := strconv.Atoi(opStr)
		if err != nil {
			return fmt.Errorf("invalid operation index: %v", err)
		}
		if op < 1 || op > t.n-2 {
			return fmt.Errorf("operation index out of range: %d", op)
		}
		idx := op - 1
		x := arr[idx] ^ arr[idx+1] ^ arr[idx+2]
		arr[idx], arr[idx+1], arr[idx+2] = x, x, x
	}
	for i, v := range arr {
		if v != 0 {
			return fmt.Errorf("array not zero at position %d", i+1)
		}
	}
	return nil
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		input := buildInputB(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nOutput:%s", i+1, err, out)
			os.Exit(1)
		}
		if err := verifyOutputB(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\noutput:%sinput:%s", i+1, err, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
