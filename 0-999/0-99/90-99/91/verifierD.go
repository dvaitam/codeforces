package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesDRaw = `5
4 1 3 2 5

7
2 6 1 3 4 7 5

4
2 4 3 1

6
1 4 5 3 6 2

2
1 2

6
6 2 4 1 5 3

1
1

1
1

2
1 2

4
1 3 4 2

5
4 2 5 3 1

5
5 2 1 3 4

7
2 4 7 5 3 1 6

5
2 5 1 4 3

7
1 7 4 3 5 6 2

7
5 3 2 4 1 6 7

3
3 2 1

7
5 3 6 2 7 4 1

3
2 1 3

4
3 2 1 4

2
1 2

7
6 4 7 5 2 3 1

7
1 3 6 2 4 5 7

2
1 2

5
5 3 4 1 2

1
1

2
1 2

6
4 3 5 1 2 6

3
1 3 2

4
4 3 2 1

4
1 4 3 2

1
1

5
1 2 3 4 5

2
1 2

3
2 3 1

6
5 4 1 2 3 6

6
6 1 5 3 2 4

7
5 3 7 2 6 1 4

1
1

1
1

6
3 6 4 1 2 5

6
2 3 6 1 5 4

7
4 3 2 1 6 7 5

2
2 1

4
2 3 4 1

5
3 1 4 5 2

5
5 2 4 1 3

7
6 2 4 3 7 1 5

5
4 5 1 3 2

3
3 2 1

3
3 1 2

3
2 1 3

3
2 3 1

7
3 7 2 1 4 5 6

7
2 6 1 3 7 5 4

6
2 6 3 1 4 5

6
5 3 4 2 6 1

5
4 1 3 5 2

1
1

7
5 2 7 3 6 4 1

1
1

6
1 5 6 2 4 3

5
4 3 5 1 2

5
5 3 4 2 1

6
4 3 1 6 5 2

2
2 1

3
1 3 2

6
1 4 6 5 2 3

3
3 2 1

7
3 2 6 4 5 1 7

7
1 5 3 7 6 4 2

5
2 4 5 1 3

2
1 2

4
4 1 3 2

5
3 4 5 2 1

1
1

3
2 3 1

5
1 4 3 2 5

3
3 2 1

2
2 1

5
1 5 3 4 2

7
7 5 6 3 2 1 4

6
6 5 2 4 3 1

1
1

3
2 1 3

3
3 1 2

3
2 3 1

7
5 6 3 1 4 7 2

6
6 5 1 4 3 2

5
3 5 2 1 4

5
5 2 4 1 3

6
4 6 1 2 3 5

3
1 2 3

2
1 2

3
2 3 1

1
1

4
1 3 4 2

3
3 1 2

7
5 6 1 7 3 4 2

4
1 3 2 4

`

// ---------- embedded solver (from cf_t25_91_D.go) ----------

func solveD91(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)

	readInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	if !scanner.Scan() {
		return ""
	}
	n, _ := strconv.Atoi(scanner.Text())

	A := make([]int, n+1)
	for i := 1; i <= n; i++ {
		A[i] = readInt()
	}

	var opsB [][]int
	var opsC [][]int

	addOp := func(cycles [][]int) {
		var b, c []int
		for _, cyc := range cycles {
			b = append(b, cyc...)
			c = append(c, cyc[1:]...)
			c = append(c, cyc[0])
		}
		opsB = append(opsB, b)
		opsC = append(opsC, c)
	}

	executeOp := func(b, c []int) {
		opsB = append(opsB, b)
		opsC = append(opsC, c)
	}

	var cycles2 [][]int
	var cycles3 [][]int
	var cycles4 [][]int

	visited := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if !visited[i] {
			curr := i
			var cyc []int
			for !visited[curr] {
				visited[curr] = true
				cyc = append(cyc, curr)
				curr = A[curr]
			}

			if len(cyc) >= 2 {
				idx := 0
				for len(cyc)-idx >= 6 {
					b := []int{cyc[idx], cyc[idx+1], cyc[idx+2], cyc[idx+3], cyc[idx+4]}
					c := []int{b[1], b[2], b[3], b[4], b[0]}
					executeOp(b, c)
					cyc[idx+4] = cyc[idx]
					idx += 4
				}
				remCyc := cyc[idx:]
				if len(remCyc) == 5 {
					addOp([][]int{remCyc})
				} else if len(remCyc) == 4 {
					cycles4 = append(cycles4, remCyc)
				} else if len(remCyc) == 3 {
					cycles3 = append(cycles3, remCyc)
				} else if len(remCyc) == 2 {
					cycles2 = append(cycles2, remCyc)
				}
			}
		}
	}

	for {
		if len(cycles2) >= 1 && len(cycles3) >= 1 {
			u := cycles2[len(cycles2)-1]
			cycles2 = cycles2[:len(cycles2)-1]
			v := cycles3[len(cycles3)-1]
			cycles3 = cycles3[:len(cycles3)-1]
			addOp([][]int{u, v})
		} else if len(cycles3) >= 2 {
			u := cycles3[len(cycles3)-1]
			cycles3 = cycles3[:len(cycles3)-1]
			v := cycles3[len(cycles3)-1]
			cycles3 = cycles3[:len(cycles3)-1]

			b := append([]int{}, u...)
			b = append(b, v[0], v[1])

			c := append([]int{}, u[1:]...)
			c = append(c, u[0])
			c = append(c, v[1], v[0])

			executeOp(b, c)
			cycles2 = append(cycles2, []int{v[0], v[2]})
		} else if len(cycles2) >= 2 {
			u := cycles2[len(cycles2)-1]
			cycles2 = cycles2[:len(cycles2)-1]
			v := cycles2[len(cycles2)-1]
			cycles2 = cycles2[:len(cycles2)-1]
			addOp([][]int{u, v})
		} else if len(cycles4) >= 1 {
			u := cycles4[len(cycles4)-1]
			cycles4 = cycles4[:len(cycles4)-1]
			addOp([][]int{u})
		} else if len(cycles3) >= 1 {
			u := cycles3[len(cycles3)-1]
			cycles3 = cycles3[:len(cycles3)-1]
			addOp([][]int{u})
		} else if len(cycles2) >= 1 {
			u := cycles2[len(cycles2)-1]
			cycles2 = cycles2[:len(cycles2)-1]
			addOp([][]int{u})
		} else {
			break
		}
	}

	var out bytes.Buffer
	fmt.Fprintln(&out, len(opsB))
	for i := 0; i < len(opsB); i++ {
		fmt.Fprintln(&out, len(opsB[i]))
		for j, x := range opsB[i] {
			if j > 0 {
				fmt.Fprint(&out, " ")
			}
			fmt.Fprint(&out, x)
		}
		fmt.Fprintln(&out)
		for j, x := range opsC[i] {
			if j > 0 {
				fmt.Fprint(&out, " ")
			}
			fmt.Fprint(&out, x)
		}
		fmt.Fprintln(&out)
	}
	return strings.TrimSpace(out.String())
}

// ---------- end embedded solver ----------

func runProg(bin, input string) (string, error) {
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

// verifyOutput checks that the candidate output is a valid solution:
// each operation is a permutation of its b-array, and applying all operations
// transforms identity into the target permutation A.
func verifyOutput(n int, A []int, output string) error {
	sc := bufio.NewScanner(strings.NewReader(output))
	sc.Split(bufio.ScanWords)
	rdInt := func() (int, error) {
		if !sc.Scan() {
			return 0, fmt.Errorf("unexpected end of output")
		}
		return strconv.Atoi(sc.Text())
	}

	numOps, err := rdInt()
	if err != nil {
		return fmt.Errorf("reading numOps: %v", err)
	}

	// start with identity permutation
	perm := make([]int, n+1)
	for i := 1; i <= n; i++ {
		perm[i] = i
	}

	for op := 0; op < numOps; op++ {
		k, err := rdInt()
		if err != nil {
			return fmt.Errorf("op %d: reading k: %v", op+1, err)
		}
		if k < 1 || k > 5 {
			return fmt.Errorf("op %d: k=%d out of range [1,5]", op+1, k)
		}
		b := make([]int, k)
		c := make([]int, k)
		for i := 0; i < k; i++ {
			b[i], err = rdInt()
			if err != nil {
				return fmt.Errorf("op %d: reading b[%d]: %v", op+1, i, err)
			}
		}
		for i := 0; i < k; i++ {
			c[i], err = rdInt()
			if err != nil {
				return fmt.Errorf("op %d: reading c[%d]: %v", op+1, i, err)
			}
		}
		// b and c must be permutations of each other
		bset := make(map[int]int)
		for _, v := range b {
			bset[v]++
		}
		cset := make(map[int]int)
		for _, v := range c {
			cset[v]++
		}
		for k, v := range bset {
			if cset[k] != v {
				return fmt.Errorf("op %d: b and c are not permutations of each other", op+1)
			}
		}
		// Apply: for each position that has value b[i], change it to c[i]
		// Build mapping b[i] -> c[i]
		mapping := make(map[int]int)
		for i := 0; i < k; i++ {
			mapping[b[i]] = c[i]
		}
		for i := 1; i <= n; i++ {
			if newVal, ok := mapping[perm[i]]; ok {
				perm[i] = newVal
			}
		}
	}

	// Check result matches A
	for i := 1; i <= n; i++ {
		if perm[i] != A[i] {
			return fmt.Errorf("after all ops, perm[%d]=%d but expected %d", i, perm[i], A[i])
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

	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	var lines []string
	process := func() {
		if len(lines) == 0 {
			return
		}
		idx++
		input := strings.Join(lines, "\n") + "\n"

		// Get reference number of ops
		expStr := solveD91(input)

		// Parse expected num ops
		expOps := 0
		if expStr != "" {
			sc := bufio.NewScanner(strings.NewReader(expStr))
			sc.Scan()
			expOps, _ = strconv.Atoi(strings.TrimSpace(sc.Text()))
		}

		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}

		// Parse got num ops
		gotOps := 0
		if got != "" {
			sc := bufio.NewScanner(strings.NewReader(got))
			sc.Scan()
			gotOps, _ = strconv.Atoi(strings.TrimSpace(sc.Text()))
		}

		// Number of ops must match (optimal)
		if gotOps != expOps {
			fmt.Printf("case %d: expected %d ops got %d ops\n", idx, expOps, gotOps)
			os.Exit(1)
		}

		// Parse n and A from input
		sc := bufio.NewScanner(strings.NewReader(input))
		sc.Split(bufio.ScanWords)
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		A := make([]int, n+1)
		for i := 1; i <= n; i++ {
			sc.Scan()
			A[i], _ = strconv.Atoi(sc.Text())
		}

		// Verify the candidate output is correct
		if err := verifyOutput(n, A, got); err != nil {
			fmt.Printf("case %d verification failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			process()
			lines = nil
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	process()
	fmt.Printf("All %d tests passed\n", idx)
}
