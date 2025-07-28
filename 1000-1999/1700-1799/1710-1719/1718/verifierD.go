package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func genCase(rng *rand.Rand) string {
	t := 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	n := rng.Intn(4) + 2
	q := rng.Intn(3) + 1
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	// permutation p
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", perm[i]+1))
	}
	sb.WriteByte('\n')
	// array a with at least 2 zeros
	k := rng.Intn(n-1) + 2
	zeros := rand.Perm(n)[:k]
	arr := make([]int, n)
	used := map[int]bool{}
	for i := 0; i < n; i++ {
		arr[i] = 0
	}
	for i := 0; i < n; i++ {
		if !contains(zeros, i) {
			val := rng.Intn(15) + 1
			for used[val] {
				val = rng.Intn(15) + 1
			}
			used[val] = true
			arr[i] = val
		}
	}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	// set S
	numbers := []int{}
	for len(numbers) < k-1 {
		v := rng.Intn(15) + 1
		if !used[v] {
			used[v] = true
			numbers = append(numbers, v)
		}
	}
	for i, v := range numbers {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		d := rng.Intn(15) + 1
		sb.WriteString(fmt.Sprintf("%d\n", d))
	}
	return sb.String()
}

func contains(arr []int, x int) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func isSimilar(p, arr []int) bool {
	n := len(arr)
	for l := 0; l < n; l++ {
		maxpIdx := l
		maxp := p[l]
		for r := l; r < n; r++ {
			if p[r] > maxp {
				maxp = p[r]
				maxpIdx = r
			}
			maxaIdx := l
			maxa := arr[l]
			for i := l; i <= r; i++ {
				if arr[i] > maxa {
					maxa = arr[i]
					maxaIdx = i
				}
			}
			if maxaIdx != maxpIdx {
				return false
			}
		}
	}
	return true
}

func brute(input string) (string, error) {
	var t int
	fmt.Sscan(input, &t)
	lines := strings.Split(strings.TrimSpace(input), "\n")
	idx := 1
	var outputs []string
	for caseNum := 0; caseNum < t; caseNum++ {
		var n, q int
		fmt.Sscanf(lines[idx], "%d %d", &n, &q)
		idx++
		p := make([]int, n)
		vals := strings.Fields(lines[idx])
		for i := 0; i < n; i++ {
			fmt.Sscanf(vals[i], "%d", &p[i])
		}
		idx++
		a := make([]int, n)
		vals = strings.Fields(lines[idx])
		for i := 0; i < n; i++ {
			fmt.Sscanf(vals[i], "%d", &a[i])
		}
		idx++
		Svals := strings.Fields(lines[idx])
		S := make([]int, len(Svals))
		for i := range Svals {
			fmt.Sscanf(Svals[i], "%d", &S[i])
		}
		idx++
		zeros := []int{}
		used := map[int]bool{}
		for i, v := range a {
			if v == 0 {
				zeros = append(zeros, i)
			} else {
				used[v] = true
			}
		}
		for ; q > 0; q-- {
			dLine := lines[idx]
			idx++
			var d int
			fmt.Sscanf(dLine, "%d", &d)
			nums := append([]int{}, S...)
			nums = append(nums, d)
			if len(nums) != len(zeros) {
				outputs = append(outputs, "NO")
				continue
			}
			usedPerm := make([]bool, len(nums))
			arr := make([]int, len(a))
			copy(arr, a)
			found := false
			var dfs func(int)
			dfs = func(pos int) {
				if found {
					return
				}
				if pos == len(zeros) {
					uniq := map[int]bool{}
					for _, v := range arr {
						if uniq[v] {
							return
						}
						uniq[v] = true
					}
					if isSimilar(p, arr) {
						found = true
					}
					return
				}
				for i, v := range nums {
					if usedPerm[i] {
						continue
					}
					arr[zeros[pos]] = v
					usedPerm[i] = true
					dfs(pos + 1)
					usedPerm[i] = false
				}
			}
			dfs(0)
			if found {
				outputs = append(outputs, "YES")
			} else {
				outputs = append(outputs, "NO")
			}
		}
	}
	return strings.Join(outputs, "\n"), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := brute(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "brute error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
