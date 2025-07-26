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

func expectedSequence(n int, a, b []int) (bool, []int) {
	dpPrev := [4]bool{}
	parent := make([][4]int, n)
	for v := 0; v < 4; v++ {
		dpPrev[v] = true
		parent[0][v] = -1
	}
	for i := 0; i < n-1; i++ {
		var dpCurr [4]bool
		for v := 0; v < 4; v++ {
			if !dpPrev[v] {
				continue
			}
			for w := 0; w < 4; w++ {
				if (v|w) == a[i] && (v&w) == b[i] {
					if !dpCurr[w] {
						dpCurr[w] = true
						parent[i+1][w] = v
					}
				}
			}
		}
		dpPrev = dpCurr
	}
	endVal := -1
	for v := 0; v < 4; v++ {
		if dpPrev[v] {
			endVal = v
			break
		}
	}
	if endVal < 0 {
		return false, nil
	}
	t := make([]int, n)
	t[n-1] = endVal
	for i := n - 1; i > 0; i-- {
		t[i-1] = parent[i][t[i]]
	}
	return true, t
}

func generateCase(rng *rand.Rand) (string, bool) {
	n := rng.Intn(18) + 2 // 2..19
	a := make([]int, n-1)
	b := make([]int, n-1)
	for i := 0; i < n-1; i++ {
		a[i] = rng.Intn(4)
		b[i] = rng.Intn(4)
	}
	exist, _ := expectedSequence(n, a, b)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n-1; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n-1; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(b[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), exist
}

func parseSequence(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) < 1 {
		return nil, fmt.Errorf("empty output")
	}
	if strings.ToUpper(fields[0]) == "NO" {
		if len(fields) != 1 {
			return nil, fmt.Errorf("unexpected extra tokens after NO")
		}
		return nil, nil
	}
	if strings.ToUpper(fields[0]) != "YES" {
		return nil, fmt.Errorf("expected YES/NO got %q", fields[0])
	}
	if len(fields) != 1+n {
		return nil, fmt.Errorf("expected %d numbers got %d", n, len(fields)-1)
	}
	seq := make([]int, n)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[1+i])
		if err != nil {
			return nil, fmt.Errorf("invalid int %q", fields[1+i])
		}
		if v < 0 || v > 3 {
			return nil, fmt.Errorf("value out of range: %d", v)
		}
		seq[i] = v
	}
	return seq, nil
}

func runCase(bin, input string, expectExist bool, a, b []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errBuf.String())
	}
	seq, err := parseSequence(out.String(), len(b)+1)
	if err != nil {
		return err
	}
	if seq == nil {
		if expectExist {
			return fmt.Errorf("expected YES but got NO")
		}
		return nil
	}
	if !expectExist {
		return fmt.Errorf("expected NO but got YES")
	}
	// verify sequence
	n := len(b) + 1
	for i := 0; i < n-1; i++ {
		if (seq[i]|seq[i+1]) != a[i] || (seq[i]&seq[i+1]) != b[i] {
			return fmt.Errorf("sequence invalid at %d", i)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		input, exist := generateCase(rng)
		// parse arrays for validation
		lines := strings.Split(strings.TrimSpace(input), "\n")
		n, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
		aStr := strings.Fields(lines[1])
		bStr := strings.Fields(lines[2])
		aArr := make([]int, n-1)
		bArr := make([]int, n-1)
		for j := 0; j < n-1; j++ {
			aArr[j], _ = strconv.Atoi(aStr[j])
			bArr[j], _ = strconv.Atoi(bStr[j])
		}
		if err := runCase(bin, input, exist, aArr, bArr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
