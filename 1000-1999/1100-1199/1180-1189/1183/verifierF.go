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

func solve(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))

	readInt := func() int {
		n := 0
		b, err := reader.ReadByte()
		if err != nil {
			return 0
		}
		for b < '0' || b > '9' {
			b, err = reader.ReadByte()
			if err != nil {
				return 0
			}
		}
		for b >= '0' && b <= '9' {
			n = n*10 + int(b-'0')
			b, err = reader.ReadByte()
			if err != nil {
				break
			}
		}
		return n
	}

	var out strings.Builder

	q := readInt()
	for t := 0; t < q; t++ {
		n := readInt()

		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = readInt()
		}

		sort.Slice(a, func(i, j int) bool {
			return a[i] > a[j]
		})

		var A []int
		for i := 0; i < n; i++ {
			if i == 0 || a[i] != a[i-1] {
				A = append(A, a[i])
			}
		}

		maxSum := 0
		m := len(A)

		for i := 0; i < m; i++ {
			x := A[i]
			valY := 0
			valZ := 0
			if i+1 < m {
				valY = A[i+1]
			}
			if i+2 < m {
				valZ = A[i+2]
			}
			if maxSum >= x+valY+valZ {
				break
			}
			if x > maxSum {
				maxSum = x
			}
			for j := i + 1; j < m; j++ {
				y := A[j]
				valZ2 := 0
				if j+1 < m {
					valZ2 = A[j+1]
				}
				if maxSum >= x+y+valZ2 {
					break
				}
				if x%y == 0 {
					continue
				}
				if x+y > maxSum {
					maxSum = x + y
				}
				for k := j + 1; k < m; k++ {
					z := A[k]
					if x%z == 0 || y%z == 0 {
						continue
					}
					if x+y+z > maxSum {
						maxSum = x + y + z
					}
					break
				}
			}
		}
		fmt.Fprintln(&out, maxSum)
	}

	return strings.TrimSpace(out.String())
}

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		n := rng.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Intn(50) + 2
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp := solve(input)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed:\nexpected: %s\ngot: %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
