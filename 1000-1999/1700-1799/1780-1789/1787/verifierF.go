package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded reference solver (same logic as cf_t24_1787_F.go)

func refPowMod(base, exp, mod int) int {
	if mod == 1 {
		return 0
	}
	res := 1
	base = base % mod
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return res
}

func refProcessChunk(chunk [][]int, Lprime, m, k int, ans []int) {
	L := Lprime * m
	C := make([]int, L)
	step := refPowMod(2, k, L)
	for j := 0; j < m; j++ {
		for x := 0; x < Lprime; x++ {
			idx := (j + x*step) % L
			C[idx] = chunk[j][x]
		}
	}
	for i := 0; i < L; i++ {
		ans[C[i]] = C[(i+1)%L]
	}
}

func refSolve(n, k int, a []int) (bool, []int) {
	visited := make([]bool, n+1)
	cyclesByLen := make(map[int][][]int)

	for i := 1; i <= n; i++ {
		if !visited[i] {
			cycle := []int{}
			curr := i
			for !visited[curr] {
				visited[curr] = true
				cycle = append(cycle, curr)
				curr = a[curr]
			}
			cyclesByLen[len(cycle)] = append(cyclesByLen[len(cycle)], cycle)
		}
	}

	ans := make([]int, n+1)
	req := 0
	if k < 30 {
		req = 1 << k
	} else {
		req = n + 1
	}

	possible := true

	for Lprime, list := range cyclesByLen {
		count := len(list)
		if Lprime%2 == 0 {
			if count%req != 0 {
				possible = false
				break
			}
			for i := 0; i < count; i += req {
				chunk := list[i : i+req]
				refProcessChunk(chunk, Lprime, req, k, ans)
			}
		} else {
			idx := 0
			if req <= count {
				numReq := count / req
				for i := 0; i < numReq; i++ {
					chunk := list[idx : idx+req]
					refProcessChunk(chunk, Lprime, req, k, ans)
					idx += req
				}
			}
			rem := count % req
			for c := 0; c < 30; c++ {
				if (rem & (1 << c)) != 0 {
					chunk := list[idx : idx+(1<<c)]
					refProcessChunk(chunk, Lprime, 1<<c, k, ans)
					idx += (1 << c)
				}
			}
		}
	}

	if !possible {
		return false, nil
	}
	return true, ans
}

// Apply permutation p 2^k times to get the result
func applyPow2k(p []int, n, k int) []int {
	// p is 1-indexed: p[i] for i in 1..n
	cur := make([]int, n+1)
	copy(cur, p)
	for step := 0; step < k; step++ {
		next := make([]int, n+1)
		for i := 1; i <= n; i++ {
			next[i] = cur[cur[i]]
		}
		cur = next
	}
	return cur
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) (string, int, int, []int) {
	n := r.Intn(7) + 1
	k := r.Intn(5) + 1
	perm := rand.Perm(n)
	for i := range perm {
		perm[i]++
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", perm[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), n, k, perm
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 200; i++ {
		input, n, k, aPrime := genCase(r)

		// Get reference answer
		aMap := make([]int, n+1)
		for j := 0; j < n; j++ {
			aMap[j+1] = aPrime[j]
		}
		refOk, _ := refSolve(n, k, aMap)

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}

		lines := strings.Split(got, "\n")
		firstLine := strings.TrimSpace(lines[0])

		if !refOk {
			if firstLine != "NO" {
				fmt.Fprintf(os.Stderr, "case %d: expected NO, got %s\ninput:\n%s", i, firstLine, input)
				os.Exit(1)
			}
			continue
		}

		// Reference says YES, candidate must say YES too
		if firstLine != "YES" {
			fmt.Fprintf(os.Stderr, "case %d: expected YES, got %s\ninput:\n%s", i, firstLine, input)
			os.Exit(1)
		}

		if len(lines) < 2 {
			fmt.Fprintf(os.Stderr, "case %d: missing permutation line\ninput:\n%s", i, input)
			os.Exit(1)
		}

		// Parse candidate permutation
		parts := strings.Fields(lines[1])
		if len(parts) != n {
			fmt.Fprintf(os.Stderr, "case %d: expected %d elements, got %d\ninput:\n%s", i, n, len(parts), input)
			os.Exit(1)
		}

		candPerm := make([]int, n+1)
		used := make([]bool, n+1)
		valid := true
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(parts[j])
			if err != nil || v < 1 || v > n || used[v] {
				valid = false
				break
			}
			used[v] = true
			candPerm[j+1] = v
		}
		if !valid {
			fmt.Fprintf(os.Stderr, "case %d: invalid permutation\ninput:\n%sgot: %s\n", i, input, lines[1])
			os.Exit(1)
		}

		// Verify: candPerm^{2^k} == aPrime
		result := applyPow2k(candPerm, n, k)
		for j := 1; j <= n; j++ {
			if result[j] != aMap[j] {
				fmt.Fprintf(os.Stderr, "case %d: a^(2^k) != a' at position %d: got %d, expected %d\ninput:\n%s", i, j, result[j], aMap[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All 200 tests passed")
}
