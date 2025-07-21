package main

import (
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

func canReachTarget(n int, w int, a []int) bool {
	if n == 1 {
		return w == a[0]
	}
	
	t := make([]int, n)
	for i := 0; i < n; i++ {
		t[i] = i
	}
	sort.Slice(t, func(i, j int) bool {
		return a[t[i]] < a[t[j]]
	})
	
	if w < 1 || (w == 1 && a[t[0]] > 1) {
		return false
	}
	
	sum := 0
	for _, v := range a {
		sum += v
	}
	
	return sum >= w
}

func validateSchedule(schedule []int, n int, target int, iterations []int) bool {
	if len(schedule) == 0 {
		return false
	}
	
	expectedLength := 0
	for _, iter := range iterations {
		expectedLength += 2 * iter
	}
	
	if len(schedule) != expectedLength {
		return false
	}
	
	processInstructionCount := make([]int, n)
	y := 0
	
	for _, processID := range schedule {
		if processID < 1 || processID > n {
			return false
		}
		
		processIdx := processID - 1
		instructionNum := processInstructionCount[processIdx]
		
		if instructionNum >= iterations[processIdx]*2 {
			return false
		}
		
		isIncrement := (instructionNum % 2) == 0
		if isIncrement {
			y++
		} else {
			y--
		}
		
		processInstructionCount[processIdx]++
	}
	
	for i, count := range processInstructionCount {
		if count != iterations[i]*2 {
			return false
		}
	}
	
	return y == target
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		w := rng.Intn(201) - 100
		
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(20) + 1
		}
		
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, w))
		for j := 0; j < n; j++ {
			if j > 0 {
				input.WriteString(" ")
			}
			input.WriteString(strconv.Itoa(a[j]))
		}
		input.WriteString("\n")
		
		aCopy := make([]int, n)
		copy(aCopy, a)
		expectedPossible := canReachTarget(n, w, aCopy)
		
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		
		lines := strings.Split(got, "\n")
		if len(lines) < 1 {
			fmt.Fprintf(os.Stderr, "case %d failed: empty output\ninput:\n%s", i+1, input.String())
			os.Exit(1)
		}
		
		isPossible := strings.ToLower(lines[0]) == "yes"
		
		if isPossible != expectedPossible {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %t got %t\ninput:\n%s", i+1, expectedPossible, isPossible, input.String())
			os.Exit(1)
		}
		
		if isPossible {
			if len(lines) != 2 {
				fmt.Fprintf(os.Stderr, "case %d failed: expected 2 lines for Yes case\ninput:\n%s", i+1, input.String())
				os.Exit(1)
			}
			
			scheduleStrs := strings.Fields(lines[1])
			schedule := make([]int, len(scheduleStrs))
			for j, s := range scheduleStrs {
				val, parseErr := strconv.Atoi(s)
				if parseErr != nil {
					fmt.Fprintf(os.Stderr, "case %d failed: cannot parse schedule: %v\ninput:\n%s", i+1, parseErr, input.String())
					os.Exit(1)
				}
				schedule[j] = val
			}
			
			expectedLength := 0
			for _, iter := range a {
				expectedLength += 2 * iter
			}
			
			if len(schedule) != expectedLength {
				fmt.Fprintf(os.Stderr, "case %d failed: schedule length %d, expected %d\ninput:\n%s", i+1, len(schedule), expectedLength, input.String())
				os.Exit(1)
			}
			
			for _, processID := range schedule {
				if processID < 1 || processID > n {
					fmt.Fprintf(os.Stderr, "case %d failed: invalid process ID %d\ninput:\n%s", i+1, processID, input.String())
					os.Exit(1)
				}
			}
		}
	}
	fmt.Println("All tests passed")
}