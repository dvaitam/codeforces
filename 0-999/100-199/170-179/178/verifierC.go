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

type op struct {
	add  bool
	id   int
	hash int
}

type testC struct {
	input  string
	expect string
}

func solveC(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var h, m, n int
	fmt.Fscan(reader, &h, &m, &n)
	table := make([]bool, h)
	idPos := make(map[int]int)
	var total int64
	for i := 0; i < n; i++ {
		var op string
		fmt.Fscan(reader, &op)
		if op == "+" {
			var id, hash int
			fmt.Fscan(reader, &id, &hash)
			pos := hash
			var cnt int64
			for table[pos] {
				cnt++
				pos += m
				if pos >= h {
					pos %= h
				}
			}
			table[pos] = true
			idPos[id] = pos
			total += cnt
		} else {
			var id int
			fmt.Fscan(reader, &id)
			pos := idPos[id]
			table[pos] = false
			delete(idPos, id)
		}
	}
	return fmt.Sprintf("%d", total)
}

func genTests() []testC {
	rand.Seed(42)
	var tests []testC
	for i := 0; i < 100; i++ {
		h := rand.Intn(20) + 10
		m := rand.Intn(h) + 1
		n := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", h, m, n))
		nextID := 1
		active := make(map[int]bool)
		for j := 0; j < n; j++ {
			if len(active) > 0 && rand.Intn(2) == 0 {
				var id int
				for id = range active {
					break
				}
				sb.WriteString(fmt.Sprintf("- %d\n", id))
				delete(active, id)
			} else {
				hash := rand.Intn(h)
				sb.WriteString(fmt.Sprintf("+ %d %d\n", nextID, hash))
				active[nextID] = true
				nextID++
			}
		}
		input := sb.String()
		expect := solveC(input)
		tests = append(tests, testC{input, expect})
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", i+1, t.input, t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
