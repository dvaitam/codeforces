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

// Embedded correct solver for 932/D
// Problem: tree with binary lifting. Operations: add node, query max hops with budget.
func solveD(input string) string {
	data := []byte(input)
	idx := 0
	nextInt := func() int {
		n := len(data)
		for idx < n && data[idx] <= ' ' {
			idx++
		}
		if idx >= n {
			return 0
		}
		res := 0
		for idx < n && data[idx] > ' ' {
			res = res*10 + int(data[idx]-'0')
			idx++
		}
		return res
	}
	nextInt64 := func() int64 {
		n := len(data)
		for idx < n && data[idx] <= ' ' {
			idx++
		}
		if idx >= n {
			return 0
		}
		var res int64
		for idx < n && data[idx] > ' ' {
			res = res*10 + int64(data[idx]-'0')
			idx++
		}
		return res
	}

	q := nextInt()
	if q == 0 {
		return ""
	}

	jump := make([][20]int, q+5)
	sum := make([][20]int64, q+5)
	w := make([]int64, q+5)
	depth := make([]int, q+5)

	cnt := 1
	w[1] = 0
	depth[1] = 1

	var last int64

	var out strings.Builder

	for k := 0; k < q; k++ {
		type_ := nextInt()
		p := nextInt64()
		qVal := nextInt64()

		if type_ == 1 {
			R := int(p ^ last)
			W := qVal ^ last

			cnt++
			w[cnt] = W

			x := R
			if w[x] >= W {
				jump[cnt][0] = x
			} else {
				for i := 19; i >= 0; i-- {
					if depth[x] > (1<<i) && w[jump[x][i]] < W {
						x = jump[x][i]
					}
				}
				jump[cnt][0] = jump[x][0]
			}

			depth[cnt] = depth[jump[cnt][0]] + 1
			sum[cnt][0] = W

			for i := 1; i <= 19; i++ {
				jump[cnt][i] = jump[jump[cnt][i-1]][i-1]
				sum[cnt][i] = sum[cnt][i-1] + sum[jump[cnt][i-1]][i-1]
			}
		} else {
			R := int(p ^ last)
			X := qVal ^ last

			ans := 0
			curr := R

			for i := 19; i >= 0; i-- {
				if depth[curr] >= (1<<i) && sum[curr][i] <= X {
					X -= sum[curr][i]
					ans += 1 << i
					curr = jump[curr][i]
				}
			}

			last = int64(ans)
			out.WriteString(strconv.Itoa(ans))
			out.WriteByte('\n')
		}
	}

	return strings.TrimSpace(out.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCases(rng *rand.Rand) []string {
	cases := make([]string, 0, 100)

	// Generate test cases inline
	for len(cases) < 100 {
		q := rng.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", q))

		cnt := 1
		// Track which type-2 queries are valid (need at least 1 node)
		for i := 0; i < q; i++ {
			if cnt == 1 && rng.Intn(3) != 0 {
				// Bias toward type 1 early to build the tree
				R := 1
				W := int64(rng.Intn(100))
				fmt.Fprintf(&sb, "1 %d %d\n", R, W)
				cnt++
			} else if rng.Intn(2) == 0 && cnt > 1 {
				// Type 2 query
				R := rng.Intn(cnt) + 1
				X := int64(rng.Intn(200))
				fmt.Fprintf(&sb, "2 %d %d\n", R, X)
			} else {
				// Type 1 (add node)
				R := rng.Intn(cnt) + 1
				W := int64(rng.Intn(100))
				fmt.Fprintf(&sb, "1 %d %d\n", R, W)
				cnt++
			}
		}
		cases = append(cases, sb.String())
	}

	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := generateCases(rng)

	testIdx := 0
	for _, input := range cases {
		testIdx++
		expect := solveD(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s", testIdx, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\ninput:\n%s", testIdx, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testIdx)
}
