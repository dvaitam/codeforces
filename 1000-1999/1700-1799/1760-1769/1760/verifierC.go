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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("python3", bin)
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(9) + 2 // ensure n >= 2 to avoid degenerate second max
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(200) - 100
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expectForInput(input string) (string, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("malformed input")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return "", err
	}
	idx := 1
	var out strings.Builder
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx >= len(lines) {
			return "", fmt.Errorf("malformed input at case %d", caseIdx+1)
		}
		n, err := strconv.Atoi(strings.TrimSpace(lines[idx]))
		if err != nil {
			return "", err
		}
		idx++
		if idx >= len(lines) {
			return "", fmt.Errorf("missing array line for case %d", caseIdx+1)
		}
		parts := strings.Fields(lines[idx])
		if len(parts) != n {
			return "", fmt.Errorf("wrong number of elements for case %d", caseIdx+1)
		}
		idx++
		a := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(parts[i])
			if err != nil {
				return "", err
			}
			a[i] = v
		}
		// compute fm and max2 and count
		fm := a[0]
		for _, v := range a {
			if v > fm {
				fm = v
			}
		}
		countMax := 0
		for _, v := range a {
			if v == fm {
				countMax++
			}
		}
		max2 := -1 << 60
		for _, v := range a {
			if v != fm && v > max2 {
				max2 = v
			}
		}
		// In case all are equal, keep max2 as fm to produce zero differences for all non-unique max case
		if countMax == n {
			max2 = fm
		}
		// produce output line
		bw := bufio.NewWriter(&out)
		for i, v := range a {
			diff := 0
			if countMax >= 2 {
				diff = v - fm
			} else {
				if v == fm {
					diff = v - max2
				} else {
					diff = v - fm
				}
			}
			if i > 0 {
				bw.WriteByte(' ')
			}
			fmt.Fprint(bw, diff)
		}
		bw.WriteByte('\n')
		bw.Flush()
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := expectForInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "verifier expect error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
