package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveA(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(reader, &n)
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	good := make([]int, 0)
	for i := 0; i < n; i++ {
		ok := true
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if a[i][j] == 1 || a[i][j] == 3 {
				ok = false
				break
			}
		}
		if ok {
			good = append(good, i+1)
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(good)))
	if len(good) > 0 {
		sb.WriteByte('\n')
		for idx, v := range good {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
	} else {
		sb.WriteByte('\n')
	}
	return sb.String()
}

func genTests() []string {
	rand.Seed(1)
	tests := make([]string, 0, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(8) + 2 // between 2 and 9
		matrix := make([][]int, n)
		for i := 0; i < n; i++ {
			matrix[i] = make([]int, n)
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					matrix[i][j] = -1
				} else if j < i {
					// ensure symmetry
					v := matrix[j][i]
					switch v {
					case 1:
						matrix[i][j] = 2
					case 2:
						matrix[i][j] = 1
					default:
						matrix[i][j] = v
					}
				} else {
					r := rand.Intn(4) // 0..3
					matrix[i][j] = r
					switch r {
					case 0:
						matrix[j][i] = 0
					case 1:
						matrix[j][i] = 2
					case 2:
						matrix[j][i] = 1
					case 3:
						matrix[j][i] = 3
					}
				}
			}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintln(n))
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(matrix[i][j]))
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, strings.TrimSpace(sb.String()))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierA <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		expected := solveA(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
