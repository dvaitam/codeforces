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

// Correct solver for 1819D embedded directly.
func solve1819D(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024*10)

	scanWord := func() string {
		scanner.Scan()
		return scanner.Text()
	}
	scanInt := func() int {
		v, _ := strconv.Atoi(scanWord())
		return v
	}

	t := scanInt()

	const INF int = 1e9
	next_seen := make([]int, 200005)
	for i := range next_seen {
		next_seen[i] = INF
	}

	var outBuf bytes.Buffer

	for tc := 0; tc < t; tc++ {
		n := scanInt()
		m := scanInt()

		T := make([][]int, n+1)
		suffix_size := make([]int, n+2)

		for i := 1; i <= n; i++ {
			k := scanInt()
			T[i] = make([]int, k)
			for j := 0; j < k; j++ {
				T[i][j] = scanInt()
			}
		}

		for i := n; i >= 1; i-- {
			suffix_size[i] = suffix_size[i+1] + len(T[i])
		}

		nxt_unknown := make([]int, n+2)
		curr_unk := INF
		for i := n; i >= 0; i-- {
			nxt_unknown[i] = curr_unk
			if i > 0 && len(T[i]) == 0 {
				curr_unk = i
			}
		}

		nxt_conflict := make([]int, n+2)
		min_R := INF
		for i := n; i >= 1; i-- {
			R_i := INF
			for _, e := range T[i] {
				if next_seen[e] < R_i {
					R_i = next_seen[e]
				}
				next_seen[e] = i
			}
			if R_i < min_R {
				min_R = R_i
			}
			nxt_conflict[i-1] = min_R
		}

		dp := make([]int, n+2)
		dp[n] = 0

		for i := n - 1; i >= 0; i-- {
			unk := nxt_unknown[i]
			conf := nxt_conflict[i]

			if unk == INF {
				if conf == INF {
					dp[i] = suffix_size[i+1]
				} else {
					dp[i] = dp[conf]
				}
			} else {
				if conf < unk {
					dp[i] = dp[conf]
				} else {
					if conf == INF {
						dp[i] = m
					} else {
						dp[i] = dp[unk]
					}
				}
			}
		}

		outBuf.WriteString(strconv.Itoa(dp[0]))
		outBuf.WriteByte('\n')

		for i := 1; i <= n; i++ {
			for _, e := range T[i] {
				next_seen[e] = INF
			}
		}
	}

	return strings.TrimSpace(outBuf.String())
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			sb.WriteString("0\n")
			continue
		}
		k := rng.Intn(m) + 1
		arr := rng.Perm(m)[:k]
		sort.Ints(arr)
		sb.WriteString(fmt.Sprintf("%d", k))
		for _, v := range arr {
			sb.WriteString(fmt.Sprintf(" %d", v+1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := fmt.Sprintf("/tmp/%s", tag)
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candD")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp := solve1819D(input)
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
