package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type pair struct {
	a, b int64
}

type coord struct {
	x, y int64
}

type change struct {
	kind       int
	pos        coord
	prevVal    byte
	hadVal     bool
	prevRobotX int64
	prevRobotY int64
}

const (
	changeNoop = iota
	changeMove
	changeCell
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("failed to read input: %v\n", err)
		os.Exit(1)
	}
	tests, err := parseInput(inputData)
	if err != nil {
		fmt.Printf("failed to parse input: %v\n", err)
		os.Exit(1)
	}

	out, err := run(bin, string(inputData))
	if err != nil {
		fmt.Printf("submission runtime error: %v\n", err)
		os.Exit(1)
	}

	program := strings.TrimSpace(out)
	if err := validateProgram(program); err != nil {
		fmt.Printf("invalid program: %v\n", err)
		os.Exit(1)
	}

	for idx, t := range tests {
		if err := simulate(program, t.a, t.b); err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func parseInput(data []byte) ([]pair, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]pair, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &tests[i].a, &tests[i].b); err != nil {
			return nil, err
		}
	}
	return tests, nil
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

func validateProgram(program string) error {
	if len(program) > 100000 {
		return fmt.Errorf("program length %d exceeds 100000", len(program))
	}
	allowed := "01eslrudt"
	for i := 0; i < len(program); i++ {
		if !strings.ContainsRune(allowed, rune(program[i])) {
			return fmt.Errorf("illegal character %q in program", program[i])
		}
	}
	return nil
}

func simulate(program string, a, b int64) error {
	grid := make(map[coord]byte)
	placeNumber(grid, a, 1)
	placeNumber(grid, b, 0)
	rx, ry := int64(0), int64(0)
	history := make([]change, 0, len(program))

	for i := 0; i < len(program); i++ {
		cmd := program[i]
		switch cmd {
		case '0', '1':
			pos := coord{rx, ry}
			prevVal, had := grid[pos]
			grid[pos] = cmd - '0'
			history = append(history, change{
				kind:    changeCell,
				pos:     pos,
				prevVal: prevVal,
				hadVal:  had,
			})
		case 'e':
			pos := coord{rx, ry}
			prevVal, had := grid[pos]
			if had {
				delete(grid, pos)
			}
			history = append(history, change{
				kind:    changeCell,
				pos:     pos,
				prevVal: prevVal,
				hadVal:  had,
			})
		case 'l':
			prevX, prevY := rx, ry
			rx--
			history = append(history, change{
				kind:       changeMove,
				prevRobotX: prevX,
				prevRobotY: prevY,
			})
		case 'r':
			prevX, prevY := rx, ry
			rx++
			history = append(history, change{
				kind:       changeMove,
				prevRobotX: prevX,
				prevRobotY: prevY,
			})
		case 'u':
			prevX, prevY := rx, ry
			ry++
			history = append(history, change{
				kind:       changeMove,
				prevRobotX: prevX,
				prevRobotY: prevY,
			})
		case 'd':
			prevX, prevY := rx, ry
			ry--
			history = append(history, change{
				kind:       changeMove,
				prevRobotX: prevX,
				prevRobotY: prevY,
			})
		case 's':
			history = append(history, change{kind: changeNoop})
		case 't':
			pos := coord{rx, ry}
			steps := 0
			if val, ok := grid[pos]; ok {
				steps = int(val) + 1
			}
			if steps > len(history) {
				steps = len(history)
			}
			for steps > 0 {
				history = undo(history, grid, &rx, &ry)
				steps--
			}
		}
	}

	pos := coord{rx, ry}
	if _, ok := grid[pos]; !ok {
		return fmt.Errorf("robot ended on empty cell")
	}
	digits := make([]byte, 0)
	x := rx
	for {
		if val, ok := grid[coord{x, ry}]; ok {
			digits = append(digits, '0'+val)
			x++
		} else {
			break
		}
	}
	sumStr := strconv.FormatInt(a+b, 2)
	if string(digits) != sumStr {
		return fmt.Errorf("expected %s but got %s", sumStr, string(digits))
	}
	return nil
}

func undo(history []change, grid map[coord]byte, rx, ry *int64) []change {
	last := history[len(history)-1]
	history = history[:len(history)-1]
	switch last.kind {
	case changeMove:
		*rx = last.prevRobotX
		*ry = last.prevRobotY
	case changeCell:
		if last.hadVal {
			grid[last.pos] = last.prevVal
		} else {
			delete(grid, last.pos)
		}
	case changeNoop:
		// nothing to do
	}
	return history
}

func placeNumber(grid map[coord]byte, value int64, row int64) {
	bits := strconv.FormatInt(value, 2)
	offset := len(bits) - 1
	for i := 0; i < len(bits); i++ {
		x := int64(i - offset)
		grid[coord{x: x, y: row}] = bits[i] - '0'
	}
}
