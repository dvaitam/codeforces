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
	"time"
)

type operation struct {
	op  byte
	val int
}

type testB struct {
	n    int
	init []int
	q    int
	ops  []operation
}

func generateTests() []testB {
	r := rand.New(rand.NewSource(43))
	tests := make([]testB, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(10) + 1
		init := make([]int, n)
		cnt := map[int]int{}
		for j := 0; j < n; j++ {
			v := r.Intn(10) + 1
			init[j] = v
			cnt[v]++
		}
		q := r.Intn(20) + 1
		ops := make([]operation, q)
		for j := 0; j < q; j++ {
			if len(cnt) == 0 || r.Intn(2) == 0 {
				v := r.Intn(10) + 1
				cnt[v]++
				ops[j] = operation{'+', v}
			} else {
				// choose a value with count>0
				keys := make([]int, 0, len(cnt))
				for k := range cnt {
					if cnt[k] > 0 {
						keys = append(keys, k)
					}
				}
				if len(keys) == 0 {
					v := r.Intn(10) + 1
					cnt[v]++
					ops[j] = operation{'+', v}
				} else {
					v := keys[r.Intn(len(keys))]
					cnt[v]--
					if cnt[v] == 0 {
						delete(cnt, v)
					}
					ops[j] = operation{'-', v}
				}
			}
		}
		tests[i] = testB{n: n, init: init, q: q, ops: ops}
	}
	return tests
}

func expected(t testB) []string {
	cnt := make(map[int]int)
	P, Q := 0, 0
	for _, x := range t.init {
		c := cnt[x]
		P -= c / 2
		Q -= c / 4
		c++
		cnt[x] = c
		P += c / 2
		Q += c / 4
	}
	res := make([]string, t.q)
	for i, op := range t.ops {
		c := cnt[op.val]
		P -= c / 2
		Q -= c / 4
		if op.op == '+' {
			c++
		} else {
			c--
		}
		cnt[op.val] = c
		P += c / 2
		Q += c / 4
		if Q >= 2 || (Q >= 1 && P >= 4) {
			res[i] = "YES"
		} else {
			res[i] = "NO"
		}
	}
	return res
}

func (t testB) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i, v := range t.init {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", t.q))
	for _, op := range t.ops {
		sb.WriteString(fmt.Sprintf("%c %d\n", op.op, op.val))
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, tc := range tests {
		expLines := expected(tc)
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(strings.NewReader(gotRaw))
		var idx int
		for idx = 0; idx < len(expLines) && scanner.Scan(); idx++ {
			got := strings.TrimSpace(scanner.Text())
			if got != expLines[idx] {
				fmt.Printf("test %d failed\ninput:\n%sexpected: %s got: %s\n", i+1, tc.Input(), expLines[idx], got)
				os.Exit(1)
			}
		}
		if idx != len(expLines) || scanner.Scan() {
			fmt.Printf("test %d: wrong number of output lines\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
