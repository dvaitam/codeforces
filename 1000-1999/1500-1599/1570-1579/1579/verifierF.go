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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveF(n, d int, arr []int) string {
	g := gcd(n, d)
	ans := 0
	impossible := false
	for start := 0; start < g && !impossible; start++ {
		cycle := []int{}
		j := start
		for {
			cycle = append(cycle, arr[j])
			j = (j + d) % n
			if j == start {
				break
			}
		}
		allOne := true
		for _, v := range cycle {
			if v == 0 {
				allOne = false
				break
			}
		}
		if allOne {
			impossible = true
			break
		}
		cur, best := 0, 0
		L := len(cycle)
		for i := 0; i < 2*L; i++ {
			if cycle[i%L] == 1 {
				cur++
				if cur > best {
					best = cur
				}
			} else {
				cur = 0
			}
		}
		if best > ans {
			ans = best
		}
	}
	if impossible {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", ans)
}

func genCaseF(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	d := rng.Intn(n) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(2)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, d))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expect := solveF(n, d, arr)
	return input, expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := genCaseF(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
