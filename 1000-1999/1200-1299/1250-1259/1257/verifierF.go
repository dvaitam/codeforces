package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	sort.Ints(a)
	j := 0
	for i := 1; i < len(a); i++ {
		if a[i] != a[j] {
			j++
			a[j] = a[i]
		}
	}
	return a[:j+1]
}

func intsToKey(a []int) string {
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func expected(arr []int) int {
	sort.Ints(arr)
	arr = uniqueInts(arr)
	n := len(arr)
	mp := make(map[string]int)
	for mask := 0; mask < (1 << 15); mask++ {
		d := make([]int, n)
		for i := 0; i < n; i++ {
			d[i] = bits.OnesCount(uint(arr[i]&0x7fff ^ mask))
		}
		base := d[0]
		for i := 0; i < n; i++ {
			d[i] -= base
		}
		mp[intsToKey(d)] = mask
	}
	for mask := 0; mask < (1 << 15); mask++ {
		d := make([]int, n)
		for i := 0; i < n; i++ {
			d[i] = 30 - bits.OnesCount(uint(arr[i]>>15^mask))
		}
		base := d[0]
		for i := 0; i < n; i++ {
			d[i] -= base
		}
		key := intsToKey(d)
		if low, ok := mp[key]; ok {
			return (mask << 15) ^ low
		}
	}
	return -1
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(4) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(1 << 20)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	exp := []string{fmt.Sprintf("%d", expected(arr))}
	return sb.String(), exp
}

func runCase(bin, input string, exp []string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		if strings.TrimSpace(line) != exp[i] {
			return fmt.Errorf("line %d expected %s got %s", i+1, exp[i], strings.TrimSpace(line))
		}
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
