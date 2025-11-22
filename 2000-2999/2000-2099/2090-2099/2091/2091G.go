package main

import (
    "bufio"
    "fmt"
    "os"
)

const inf int64 = 1 << 60

func maxPower(s int64, k int) int {
    // minCur/maxCur store the minimal and maximal reachable position for each remainder mod current step
    minCur := make([]int64, k+1)
    maxCur := make([]int64, k+1)
    for i := range minCur {
        minCur[i] = inf
        maxCur[i] = -1
    }
    minCur[0], maxCur[0] = 0, 0

    step := k
    dir := 1 // +1 means moving towards s, -1 backwards; alternates every turn

    for {
        // positions after doing at least one stroke with current power
        minEnd := make([]int64, step)
        maxEnd := make([]int64, step)
        for i := 0; i < step; i++ {
            minEnd[i] = inf
            maxEnd[i] = -1
        }

        for r := 0; r < step; r++ {
            if minCur[r] > maxCur[r] {
                continue
            }
            if dir == 1 {
                // move forward
                L := minCur[r] + int64(step) // need at least one stroke
                if L > s {
                    continue
                }
                // largest position <= s with residue r
                H := s - (s-int64(r))%int64(step)
                if L > H {
                    continue
                }
                if L < minEnd[r] {
                    minEnd[r] = L
                }
                if H > maxEnd[r] {
                    maxEnd[r] = H
                }
            } else {
                // move backward
                H := maxCur[r] - int64(step)
                if H < 0 {
                    continue
                }
                L := int64(r) // smallest non-negative with this residue
                if L > H {
                    continue
                }
                if L < minEnd[r] {
                    minEnd[r] = L
                }
                if H > maxEnd[r] {
                    maxEnd[r] = H
                }
            }
        }

        targetR := int(s % int64(step))
        if minEnd[targetR] <= s && s <= maxEnd[targetR] {
            return step
        }
        if step == 1 {
            // power cannot decrease further; with power 1 any reachable position allows finishing
            return 1
        }

        nextStep := step - 1
        minNext := make([]int64, nextStep)
        maxNext := make([]int64, nextStep)
        for i := 0; i < nextStep; i++ {
            minNext[i] = inf
            maxNext[i] = -1
        }

        // Map positions after the stroke (mod step) to remainders for the next power (step-1).
        // Since gcd(step, step-1) = 1, the residue modulo (step-1) increases by 1 each time we add 'step'.
        for r := 0; r < step; r++ {
            L := minEnd[r]
            H := maxEnd[r]
            if L > H {
                continue
            }

            t0 := (L - int64(r)) / int64(step) // minimal multiplier producing L
            t1 := (H - int64(r)) / int64(step)
            length := t1 - t0 + 1

            startRes := int((int64(r%nextStep) + t0) % int64(nextStep))

            if length <= int64(nextStep) {
                pos := L
                res := startRes
                for i := int64(0); i < length; i++ {
                    if pos < minNext[res] {
                        minNext[res] = pos
                    }
                    if pos > maxNext[res] {
                        maxNext[res] = pos
                    }
                    pos += int64(step)
                    res++
                    if res == nextStep {
                        res = 0
                    }
                }
            } else {
                // Earliest occurrences: first full cycle of residues
                pos := L
                res := startRes
                for i := 0; i < nextStep; i++ {
                    if pos < minNext[res] {
                        minNext[res] = pos
                    }
                    pos += int64(step)
                    res++
                    if res == nextStep {
                        res = 0
                    }
                }
                // Latest occurrences: last full cycle near the end
                pos = H
                res = int((int64(r%nextStep) + t1) % int64(nextStep))
                for i := 0; i < nextStep; i++ {
                    if pos > maxNext[res] {
                        maxNext[res] = pos
                    }
                    pos -= int64(step)
                    res--
                    if res < 0 {
                        res += nextStep
                    }
                }
            }
        }

        minCur, maxCur = minNext, maxNext
        step = nextStep
        dir = -dir
    }
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    if _, err := fmt.Fscan(in, &T); err != nil {
        return
    }
    for ; T > 0; T-- {
        var s int64
        var k int
        fmt.Fscan(in, &s, &k)
        fmt.Fprintln(out, maxPower(s, k))
    }
}

