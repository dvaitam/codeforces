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
	
	// Init req and has with placeholder 0
	// 'M', 'C', '-' are valid.
	req := make([]byte, n)
	has := make([]byte, n)
	
	// Generate random permutation for itinerary
	perm := rand.Perm(n)
	
	passenger := byte(0) // 0 means empty
	
	for _, idx := range perm {
		// idx is building index 0..n-1
		
		// Decide action
		// If scooter empty: Pickup or Pass
		// If scooter full: Drop, Swap, or Pass
		
		action := rand.Intn(2) // 0: Pass, 1: Interact
		if passenger != 0 {
			action = rand.Intn(3) // 0: Pass, 1: Drop, 2: Swap
		}
		
		if action == 0 { // Pass
			// Do nothing to flow (will fill match later)
			continue
		}
		
		if passenger == 0 {
			// Must be Pickup
			// Decide what to pickup
			res := byte('M')
			if rand.Intn(2) == 0 {
				res = 'C'
			}
			has[idx] = res
			passenger = res
		} else {
			// Drop or Swap
			// Drop: Req = passenger. Has = ? (will be filled). Passenger = 0.
			// Swap: Req = passenger. Has = NewRes. Passenger = NewRes.
			
			req[idx] = passenger
			
			if action == 1 { // Drop
				passenger = 0
				// has[idx] remains 0 (will be filled with - or garbage)
			} else { // Swap
				res := byte('M')
				if rand.Intn(2) == 0 {
					res = 'C'
				}
				has[idx] = res
				passenger = res
			}
		}
	}
	
	// If passenger still on scooter, we must Drop it somewhere?
	// Or just retry.
	if passenger != 0 {
		return genTest()
	}
	
	// Fill remaining
	for i := 0; i < n; i++ {
		if has[i] == 0 && req[i] == 0 {
			// Untouched or Passed
			// Make them match
			res := byte('-')
			r := rand.Intn(3)
			if r == 1 { res = 'M' }
			if r == 2 { res = 'C' }
			has[i] = res
			req[i] = res
		} else if has[i] == 0 {
			// Req set (Drop target), Has not set
			// Has can be anything that doesn't need to be saved.
			// Safest is '-'
			has[i] = '-'
		} else if req[i] == 0 {
			// Has set (Pickup source), Req not set
			// Req can be anything satisfied by Has?
			// No, Has was Picked up (gone).
			// So Req must be satisfied by what?
			// Wait, if we Picked up from i, Has[i] is gone.
			// Req[i] must be satisfied by... nothing (since scooter left).
			// So Req[i] MUST be '-' (no class).
			// UNLESS we dropped something there first?
			// "At most one DROPOFF and one PICKUP... in this order".
			// My generator simulates:
			// Pickup (only).
			// Drop (only).
			// Swap (Drop then Pick).
			
			// If we did Pickup only:
			// Has[i] = M. Picked up.
			// Building i is now empty.
			// Req[i] must be satisfied. So Req[i] must be '-'.
			req[i] = '-'
			
			// If we did Swap:
			// Has[i] = C. Req[i] = M (Passenger).
			// Drop M, Pick C.
			// End state: Building has M. Req is M. Satisfied.
			// Logic above set req[idx] = passenger. has[idx] = res.
			// So req is set. This branch (req==0) won't be hit for Swap.
		}
	}
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteByte(req[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		sb.WriteByte(has[i])
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

	position := -1
	passenger := byte('-')
	
	// State for handling swaps (Drop on full building -> Pending Pick)
	swapping := false
	var swapDrop byte

	for i, op := range ops {
		fields := strings.Fields(op)
		if len(fields) == 0 {
			return fmt.Errorf("empty instruction at #%d", i+1)
		}

		switch fields[0] {
		case "DRIVE":
			if swapping {
				return fmt.Errorf("DRIVE instruction issued while in the middle of a swap (missing PICKUP) at #%d", i+1)
			}
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
			position = x - 1
		case "PICKUP":
			if position == -1 {
				return fmt.Errorf("PICKUP before visiting any building at #%d", i+1)
			}
			
			if swapping {
				// Completing a swap
				// We drop 'swapDrop' and pick up 'buildings[position]'
				if buildings[position] == '-' {
					return fmt.Errorf("logic error: swapping on empty building at #%d", i+1)
				}
				pOld := buildings[position]
				buildings[position] = swapDrop
				passenger = pOld
				
				swapping = false
				swapDrop = '-'
			} else {
				// Normal pickup
				if passenger != '-' {
					return fmt.Errorf("PICKUP attempted with passenger already on scooter at #%d", i+1)
				}
				if buildings[position] == '-' {
					return fmt.Errorf("PICKUP at building %d with no professor present", position+1)
				}
				passenger = buildings[position]
				buildings[position] = '-'
			}
			
		case "DROPOFF":
			if position == -1 {
				return fmt.Errorf("DROPOFF before visiting any building at #%d", i+1)
			}
			if swapping {
				return fmt.Errorf("DROPOFF issued while already swapping at #%d", i+1)
			}
			if passenger == '-' {
				return fmt.Errorf("DROPOFF without passenger at #%d", i+1)
			}
			
			if buildings[position] != '-' {
				// Building occupied: Start Swap
				swapping = true
				swapDrop = passenger
				passenger = '-'
			} else {
				// Normal drop
				buildings[position] = passenger
				passenger = '-'
			}
			
		default:
			return fmt.Errorf("unknown instruction %q at #%d", fields[0], i+1)
		}
	}

	if swapping {
		return fmt.Errorf("itinerary ended in the middle of a swap")
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
