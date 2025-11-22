package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, string(out))
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe string, input []byte) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genTest() []byte {
	n := rand.Intn(6) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	chars := []byte{'M', 'C', '-'}
	for i := 0; i < n; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func verify(input []byte, output string) error {
	var n int
	var classes, profs string
	if _, err := fmt.Fscan(bytes.NewReader(input), &n, &classes, &profs); err != nil {
		return fmt.Errorf("failed to parse input: %v", err)
	}

	outScanner := bufio.NewScanner(strings.NewReader(output))
	if !outScanner.Scan() {
		return fmt.Errorf("missing number of operations in output")
	}
	var l int
	if _, err := fmt.Sscan(outScanner.Text(), &l); err != nil {
		return fmt.Errorf("failed to read number of operations: %v", err)
	}

	var ops []string
	for outScanner.Scan() {
		ops = append(ops, strings.TrimSpace(outScanner.Text()))
	}
	if err := outScanner.Err(); err != nil {
		return fmt.Errorf("reading output failed: %v", err)
	}

	if len(ops) != l {
		return fmt.Errorf("expected %d operations, got %d", l, len(ops))
	}

	buildings := []byte(profs)
	demand := []byte(classes)

	visited := make([]bool, n)
	pickupUsed := make([]bool, n)
	dropoffUsed := make([]bool, n)

	position := -1
	passenger := byte('-')

	for i, op := range ops {
		fields := strings.Fields(op)
		if len(fields) == 0 {
			return fmt.Errorf("empty instruction at #%d", i+1)
		}

		switch fields[0] {
		case "DRIVE":
			if len(fields) != 2 {
				return fmt.Errorf("invalid DRIVE format at #%d", i+1)
			}
			var x int
			if _, err := fmt.Sscan(fields[1], &x); err != nil {
				return fmt.Errorf("invalid DRIVE argument at #%d: %v", i+1, err)
			}
			if x < 1 || x > n {
				return fmt.Errorf("DRIVE target out of range at #%d", i+1)
			}
			if visited[x-1] {
				return fmt.Errorf("building %d visited multiple times", x)
			}
			visited[x-1] = true
			position = x - 1
		case "PICKUP":
			if position == -1 {
				return fmt.Errorf("PICKUP before visiting any building at #%d", i+1)
			}
			if pickupUsed[position] {
				return fmt.Errorf("multiple PICKUP at building %d", position+1)
			}
			if passenger != '-' {
				return fmt.Errorf("PICKUP attempted with passenger already on scooter at #%d", i+1)
			}
			if buildings[position] == '-' {
				return fmt.Errorf("PICKUP at building %d with no professor present", position+1)
			}
			passenger = buildings[position]
			buildings[position] = '-'
			pickupUsed[position] = true
		case "DROPOFF":
			if position == -1 {
				return fmt.Errorf("DROPOFF before visiting any building at #%d", i+1)
			}
			if dropoffUsed[position] {
				return fmt.Errorf("multiple DROPOFF at building %d", position+1)
			}
			if pickupUsed[position] {
				return fmt.Errorf("DROPOFF after PICKUP at building %d", position+1)
			}
			if passenger == '-' {
				return fmt.Errorf("DROPOFF without passenger at #%d", i+1)
			}
			dropoffUsed[position] = true
			buildings[position] = passenger
			passenger = '-'
		default:
			return fmt.Errorf("unknown instruction %q at #%d", fields[0], i+1)
		}
	}

	if passenger != '-' {
		return fmt.Errorf("itinerary ended with a passenger still on the scooter")
	}

	for i := 0; i < n; i++ {
		if demand[i] == 'M' && buildings[i] != 'M' {
			return fmt.Errorf("building %d requires M but has %c", i+1, buildings[i])
		}
		if demand[i] == 'C' && buildings[i] != 'C' {
			return fmt.Errorf("building %d requires C but has %c", i+1, buildings[i])
		}
	}

	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()
	for i := 1; i <= 100; i++ {
		in := genTest()
		got, err := runProg(exe, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i, err, got)
			os.Exit(1)
		}
		if err := verify(in, got); err != nil {
			fmt.Printf("wrong answer on test %d\ninput:\n%sgot:%s\nreason: %v\n", i, string(in), got, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
