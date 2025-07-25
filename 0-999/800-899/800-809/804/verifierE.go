package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	if n%4 == 2 || n%4 == 3 {
		return "NO"
	}
	var sb strings.Builder
	sb.WriteString("YES")
	if n == 1 {
		return sb.String()
	}
	sb.WriteByte('\n')
	ans4 := [][2]int{{2, 4}, {1, 3}, {2, 3}, {1, 4}, {1, 2}, {3, 4}}
	ans5 := [][2]int{{3, 5}, {1, 3}, {4, 5}, {2, 3}, {1, 2}, {1, 5}, {2, 5}, {3, 4}, {2, 4}, {1, 4}}
	t44 := [][2]int{{4, 5}, {1, 6}, {3, 6}, {1, 7}, {1, 5}, {3, 8}, {3, 5}, {1, 8}, {2, 6}, {4, 7}, {2, 7}, {2, 5}, {3, 7}, {4, 6}, {4, 8}, {2, 8}}
	t45 := [][2]int{{3, 7}, {1, 6}, {2, 9}, {1, 7}, {1, 5}, {3, 6}, {2, 7}, {4, 9}, {3, 9}, {3, 8}, {3, 5}, {4, 5}, {4, 6}, {2, 6}, {4, 7}, {4, 8}, {1, 8}, {1, 9}, {2, 8}, {2, 5}}
	type seg struct{ l, r int }
	var v []seg
	for i := 1; i <= n; i += 4 {
		if i+3 == n-1 {
			v = append(v, seg{i, i + 4})
			break
		}
		v = append(v, seg{i, i + 3})
	}
	for i := 0; i < len(v); i++ {
		s := v[i]
		delta := s.l - 1
		size := s.r - s.l + 1
		if size == 4 {
			for _, p := range ans4 {
				sb.WriteString(fmt.Sprintf("%d %d\n", p[0]+delta, p[1]+delta))
			}
		} else {
			for _, p := range ans5 {
				sb.WriteString(fmt.Sprintf("%d %d\n", p[0]+delta, p[1]+delta))
			}
		}
		for j := i + 1; j < len(v); j++ {
			t := v[j]
			dj := t.l - 5
			size2 := t.r - t.l + 1
			if size2 == 4 {
				for _, p := range t44 {
					sb.WriteString(fmt.Sprintf("%d %d\n", p[0]+delta, p[1]+dj))
				}
			} else {
				for _, p := range t45 {
					sb.WriteString(fmt.Sprintf("%d %d\n", p[0]+delta, p[1]+dj))
				}
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genTest() string {
	r := rand.New(rand.NewSource(rand.Int63()))
	n := r.Intn(10) + 1
	return fmt.Sprintf("%d\n", n)
}

func generateTests() []string {
	rand.Seed(5)
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		tests[i] = genTest()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveE(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed. expected %s got %s\ninput:%s", i+1, expected, got, t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
