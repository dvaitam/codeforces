package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func phi(n int) int {
	res := n
	p := 2
	x := n
	for p*p <= x {
		if x%p == 0 {
			for x%p == 0 {
				x /= p
			}
			res -= res / p
		}
		p++
	}
	if x > 1 {
		res -= res / x
	}
	return res
}

func seq(x int) []int {
	s := []int{x}
	for x > 1 {
		x = phi(x)
		s = append(s, x)
	}
	return s
}

func minChanges(arr []int) int {
    lists := make([][]int, len(arr))
    for i, v := range arr {
        lists[i] = seq(v)
    }
    // Use a large sentinel and int64 to avoid overflow
    const INF int64 = 1<<60
    best := INF
    for _, l := range lists {
        for idx, val := range l {
            cost := int64(idx)
            for _, other := range lists {
                pos := -1
                for k, vv := range other {
                    if vv == val {
                        pos = k
                        break
                    }
                }
                if pos == -1 {
                    cost = INF
                    break
                }
                cost += int64(pos)
                if cost >= best { // early prune
                    break
                }
            }
            if cost < best {
                best = cost
            }
        }
    }
    if best == INF {
        return 0
    }
    return int(best)
}

func naive(n, m int, arr []int, ops [][3]int) []int {
	res := []int{}
	for _, op := range ops {
		if op[0] == 1 {
			for i := op[1] - 1; i < op[2]; i++ {
				arr[i] = phi(arr[i])
			}
		} else {
			seg := make([]int, op[2]-op[1]+1)
			copy(seg, arr[op[1]-1:op[2]])
			res = append(res, minChanges(seg))
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierE.go path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if b, err := filepath.Abs(bin); err == nil {
		bin = b
	}

	rand.Seed(5)
	const T = 100
	for tc := 0; tc < T; tc++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rand.Intn(10) + 1
		}
		ops := make([][3]int, m)
		for i := 0; i < m; i++ {
			t := rand.Intn(2) + 1
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			ops[i] = [3]int{t, l, r}
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for _, op := range ops {
			input.WriteString(fmt.Sprintf("%d %d %d\n", op[0], op[1], op[2]))
		}

		expected := naive(n, m, append([]int(nil), arr...), ops)
		out, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Printf("test %d binary error: %v\n", tc+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(expected) {
			fmt.Printf("test %d wrong number of outputs\n", tc+1)
			os.Exit(1)
		}
		for i, exp := range expected {
			got, err := strconv.Atoi(fields[i])
			if err != nil || got != exp {
				fmt.Printf("test %d failed at output %d: expected %d got %s\n", tc+1, i+1, exp, fields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
