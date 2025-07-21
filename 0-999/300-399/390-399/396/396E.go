package main

import (
   "bufio"
   "os"
   "sort"
   "strconv"
   "strings"
   "math/big"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    line, _ := reader.ReadString('\n')
    m, _ := strconv.Atoi(strings.TrimSpace(line))
    initE := make(map[int]*big.Int)
    primesList := make([]int, 0, m)
    for i := 0; i < m; i++ {
        line, _ = reader.ReadString('\n')
        parts := strings.Fields(line)
        p, _ := strconv.Atoi(parts[0])
        a := new(big.Int)
        a.SetString(parts[1], 10)
        initE[p] = a
        primesList = append(primesList, p)
    }
    line, _ = reader.ReadString('\n')
    k := new(big.Int)
    k.SetString(strings.TrimSpace(line), 10)

    // sieve for spf
    const MAXP = 1000000
    spf := make([]int, MAXP+1)
    for i := 2; i <= MAXP; i++ {
        if spf[i] == 0 {
            for j := i; j <= MAXP; j += i {
                if spf[j] == 0 {
                    spf[j] = i
                }
            }
        }
    }

    // answer exponents
    ans := make(map[int]*big.Int)
    seen := make(map[int]bool)

    // current level primes and exponents
    curr := make([]int, 0, len(primesList))
    currE := make(map[int]*big.Int)
    for _, p := range primesList {
        curr = append(curr, p)
        currE[p] = initE[p]
        seen[p] = true
    }

    // iterate levels
    h := 0
    for len(curr) > 0 {
        // if h >= k, break
        // stop after processing level h > k
        if big.NewInt(int64(h)).Cmp(k) > 0 {
            break
        }
        // k - h
        kh := new(big.Int).Sub(k, big.NewInt(int64(h)))
        nextE := make(map[int]*big.Int)
        nextPr := make([]int, 0)
        for _, p := range curr {
            e0 := currE[p]
            // rem = max(e0 - kh, 0)
            rem := new(big.Int).Sub(e0, kh)
            if rem.Sign() > 0 {
                if _, ok := ans[p]; !ok {
                    ans[p] = new(big.Int)
                }
                ans[p].Add(ans[p], rem)
            }
            // A = min(e0, kh)
            A := new(big.Int)
            if e0.Cmp(kh) <= 0 {
                A.Set(e0)
            } else {
                A.Set(kh)
            }
            if A.Sign() <= 0 {
                continue
            }
            // factor p-1
            x := p - 1
            for x > 1 {
                f := spf[x]
                cnt := 0
                for x%f == 0 {
                    x /= f
                    cnt++
                }
                // add to nextE[f] count cnt * A
                add := new(big.Int).Mul(A, big.NewInt(int64(cnt)))
                if _, ok := nextE[f]; !ok {
                    nextE[f] = new(big.Int)
                }
                nextE[f].Add(nextE[f], add)
            }
        }
        // build next curr
        for p, e := range nextE {
            if !seen[p] {
                seen[p] = true
                nextPr = append(nextPr, p)
            }
            // if multiple sources, nextE holds sum
            currE[p] = e
        }
        curr = nextPr
        h++
    }
    // output ans
    // collect sorted primes
    outPr := make([]int, 0, len(ans))
    for p := range ans {
        outPr = append(outPr, p)
    }
    sort.Ints(outPr)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    writer.WriteString(strconv.Itoa(len(outPr)))
    writer.WriteByte('\n')
    for _, p := range outPr {
        writer.WriteString(strconv.Itoa(p))
        writer.WriteByte(' ')
        writer.WriteString(ans[p].String())
        writer.WriteByte('\n')
    }
}
