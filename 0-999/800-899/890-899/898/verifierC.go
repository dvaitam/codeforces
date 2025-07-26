package main

import (
	"bufio"
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

type entry struct {
	name string
	nums []string
}

func solveC(entries []entry) map[string][]string {
	m := make(map[string]map[string]struct{})
	for _, e := range entries {
		if _, ok := m[e.name]; !ok {
			m[e.name] = make(map[string]struct{})
		}
		for _, num := range e.nums {
			m[e.name][num] = struct{}{}
		}
	}
	res := make(map[string][]string)
	for name, set := range m {
		nums := make([]string, 0, len(set))
		for num := range set {
			nums = append(nums, num)
		}
		keep := make([]string, 0, len(nums))
		for _, x := range nums {
			drop := false
			for _, y := range nums {
				if x != y && len(y) >= len(x) && strings.HasSuffix(y, x) {
					drop = true
					break
				}
			}
			if !drop {
				keep = append(keep, x)
			}
		}
		sort.Strings(keep)
		res[name] = keep
	}
	return res
}

func genCaseC(rng *rand.Rand) (string, map[string][]string) {
	n := rng.Intn(4) + 1
	entries := make([]entry, n)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("%c%c", 'a'+rune(rng.Intn(26)), 'a'+rune(rng.Intn(26)))
		k := rng.Intn(3) + 1
		nums := make([]string, k)
		for j := 0; j < k; j++ {
			nums[j] = strconv.Itoa(rng.Intn(1000))
		}
		entries[i] = entry{name: name, nums: nums}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range entries {
		fmt.Fprintf(&sb, "%s %d", e.name, len(e.nums))
		for _, num := range e.nums {
			fmt.Fprintf(&sb, " %s", num)
		}
		sb.WriteString("\n")
	}
	return sb.String(), solveC(entries)
}

func parseOutputC(out string) (map[string][]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	if !scanner.Scan() {
		return nil, fmt.Errorf("no output")
	}
	m, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return nil, fmt.Errorf("bad first line: %v", err)
	}
	res := make(map[string][]string)
	for i := 0; i < m; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("expected %d lines", m)
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			return nil, fmt.Errorf("bad line: %q", scanner.Text())
		}
		name := fields[0]
		k, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("bad count: %v", err)
		}
		if len(fields)-2 != k {
			return nil, fmt.Errorf("wrong number of phones for %s", name)
		}
		nums := make([]string, k)
		copy(nums, fields[2:])
		sort.Strings(nums)
		res[name] = nums
	}
	return res, nil
}

func runCaseC(bin string, in string, exp map[string][]string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	got, err := parseOutputC(buf.String())
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if len(got) != len(exp) {
		return fmt.Errorf("expected %d contacts got %d", len(exp), len(got))
	}
	for name, expNums := range exp {
		gnums, ok := got[name]
		if !ok {
			return fmt.Errorf("missing friend %s", name)
		}
		if len(gnums) != len(expNums) {
			return fmt.Errorf("for %s expected %v got %v", name, expNums, gnums)
		}
		for i := range gnums {
			if gnums[i] != expNums[i] {
				return fmt.Errorf("for %s expected %v got %v", name, expNums, gnums)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseC(rng)
		if err := runCaseC(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
