package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func verify(perm []int, segs [][]int) bool {
	n := len(perm)
	counter := make(map[string]int)
	for _, s := range segs {
		key := fmt.Sprint(s)
		counter[key]++
	}
	for r := 2; r <= n; r++ {
		found := false
		for l := 1; l < r; l++ {
			sub := append([]int(nil), perm[l-1:r]...)
			sort.Ints(sub)
			key := fmt.Sprint(sub)
			if counter[key] > 0 {
				counter[key]--
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func genCase(rng *rand.Rand) (string, [][]int, int) {
	n := rng.Intn(8) + 2 // 2..9
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	segs := make([][]int, n-1)
	idx := 0
	for r := 2; r <= n; r++ {
		l := rng.Intn(r-1) + 1
		seg := append([]int(nil), perm[l-1:r]...)
		sort.Ints(seg)
		segs[idx] = seg
		idx++
	}
	rng.Shuffle(len(segs), func(i, j int) { segs[i], segs[j] = segs[j], segs[i] })
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, seg := range segs {
		sb.WriteString(fmt.Sprintf("%d", len(seg)))
		for _, v := range seg {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), segs, n
}

func runCase(bin string, input string, segs [][]int, n int) error {
	out, err := runProg(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	perm := make([]int, n)
	seen := make(map[int]bool)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fields[i])
		if err != nil {
			return fmt.Errorf("invalid number: %v", err)
		}
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range", v)
		}
		if seen[v] {
			return fmt.Errorf("duplicate value %d", v)
		}
		seen[v] = true
		perm[i] = v
	}
	if !verify(perm, segs) {
		return fmt.Errorf("output does not satisfy segments")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, segs, n := genCase(rng)
		if err := runCase(bin, input, segs, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
