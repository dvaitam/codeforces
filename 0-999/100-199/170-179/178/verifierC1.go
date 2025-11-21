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

type testCaseC1 struct {
	input  string
	expect string
}

func solveC1(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var h, m, n int
	fmt.Fscan(reader, &h, &m, &n)
	table := make([]bool, h)
	idPos := make(map[int]int, n)
	var total int64
	for i := 0; i < n; i++ {
		var op string
		fmt.Fscan(reader, &op)
		if op == "+" {
			var id, hash int
			fmt.Fscan(reader, &id, &hash)
			pos := hash
			var dummy int64
			for table[pos] {
				dummy++
				pos += m
				if pos >= h {
					pos %= h
				}
			}
			table[pos] = true
			idPos[id] = pos
			total += dummy
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

func findSlot(table []bool, start, step int) int {
	h := len(table)
	if h == 0 {
		return -1
	}
	pos := start % h
	if pos < 0 {
		pos += h
	}
	for tried := 0; tried < h; tried++ {
		if !table[pos] {
			return pos
		}
		pos += step
		pos %= h
	}
	return -1
}

func chooseHash(table []bool, step int) (int, int, bool) {
	h := len(table)
	if h == 0 {
		return 0, 0, false
	}
	occupied := make([]int, 0, h)
	freeCells := make([]int, 0, h)
	for i, used := range table {
		if used {
			occupied = append(occupied, i)
		} else {
			freeCells = append(freeCells, i)
		}
	}
	rand.Shuffle(len(occupied), func(i, j int) {
		occupied[i], occupied[j] = occupied[j], occupied[i]
	})
	rand.Shuffle(len(freeCells), func(i, j int) {
		freeCells[i], freeCells[j] = freeCells[j], freeCells[i]
	})
	candidates := append(append([]int(nil), occupied...), freeCells...)
	for _, start := range candidates {
		slot := findSlot(table, start, step)
		if slot != -1 {
			return start, slot, true
		}
	}
	return 0, 0, false
}

func pickRandomID(m map[int]int) int {
	idx := rand.Intn(len(m))
	for id := range m {
		if idx == 0 {
			return id
		}
		idx--
	}
	return 0
}

func buildTest(h, m, n int) testCaseC1 {
	table := make([]bool, h)
	idPos := make(map[int]int)
	nextID := 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", h, m, n))
	for i := 0; i < n; i++ {
		doDelete := false
		if len(idPos) == 0 {
			doDelete = false
		} else if len(idPos) == h {
			doDelete = true
		} else {
			doDelete = rand.Intn(3) == 0
		}
		if !doDelete {
			hash, slot, ok := chooseHash(table, m)
			if ok {
				id := nextID
				nextID++
				sb.WriteString(fmt.Sprintf("+ %d %d\n", id, hash))
				table[slot] = true
				idPos[id] = slot
				continue
			}
			doDelete = true
		}
		if doDelete {
			if len(idPos) == 0 {
				panic("attempted deletion on empty table")
			}
			id := pickRandomID(idPos)
			sb.WriteString(fmt.Sprintf("- %d\n", id))
			pos := idPos[id]
			table[pos] = false
			delete(idPos, id)
		}
	}
	input := sb.String()
	expect := solveC1(input)
	return testCaseC1{input: input, expect: expect}
}

func genTests() []testCaseC1 {
	rand.Seed(42)
	var tests []testCaseC1
	manual := "2 1 4\n+ 1 0\n+ 2 0\n- 1\n+ 3 0\n"
	tests = append(tests, testCaseC1{
		input:  manual,
		expect: solveC1(manual),
	})
	manual = "6 2 7\n+ 1 0\n+ 2 2\n+ 3 4\n- 2\n+ 4 2\n- 1\n+ 5 0\n"
	tests = append(tests, testCaseC1{
		input:  manual,
		expect: solveC1(manual),
	})
	for i := 0; i < 200; i++ {
		h := rand.Intn(90) + 2
		m := rand.Intn(h-1) + 1
		n := rand.Intn(200) + 1
		tests = append(tests, buildTest(h, m, n))
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
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		got, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, t.input)
			os.Exit(1)
		}
		if got != t.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpect:\n%s\nactual:\n%s\n", i+1, t.input, t.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
