package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runExe(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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
	return out.String(), nil
}

func genCase(rng *rand.Rand) (string, [][]int64) {
	n := rng.Intn(4) + 2
	m := rng.Intn(4) + 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	a := make([][]int64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			v := int64(rng.Intn(16) + 1)
			a[i][j] = v
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), a
}

func isFourthPower(x int64) bool {
	if x <= 0 {
		return false
	}
	r := int64(math.Round(math.Pow(float64(x), 0.25)))
	for d := r - 2; d <= r+2; d++ {
		if d > 0 && d*d*d*d == x {
			return true
		}
	}
	return false
}

func checkOutput(a [][]int64, out string) error {
	n := len(a)
	m := len(a[0])
	data := []int64{}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		for _, f := range fields {
			val, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer %q", f)
			}
			data = append(data, val)
		}
	}
	if len(data) != n*m {
		return fmt.Errorf("expected %d numbers, got %d", n*m, len(data))
	}
	b := make([][]int64, n)
	idx := 0
	for i := 0; i < n; i++ {
		b[i] = make([]int64, m)
		for j := 0; j < m; j++ {
			val := data[idx]
			idx++
			if val < 1 || val > 1_000_000 {
				return fmt.Errorf("value out of range at (%d,%d)", i+1, j+1)
			}
			if val%a[i][j] != 0 {
				return fmt.Errorf("not multiple at (%d,%d)", i+1, j+1)
			}
			b[i][j] = val
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i+1 < n {
				d := b[i][j] - b[i+1][j]
				if d < 0 {
					d = -d
				}
				if !isFourthPower(d) {
					return fmt.Errorf("invalid diff between (%d,%d) and (%d,%d)", i+1, j+1, i+2, j+1)
				}
			}
			if j+1 < m {
				d := b[i][j] - b[i][j+1]
				if d < 0 {
					d = -d
				}
				if !isFourthPower(d) {
					return fmt.Errorf("invalid diff between (%d,%d) and (%d,%d)", i+1, j+1, i+1, j+2)
				}
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, a := genCase(rng)
		out, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkOutput(a, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
