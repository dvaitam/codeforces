package main

import (
    "bufio"
    "fmt"
    "os"
)

const maxN = 200000
const shift = 21
const mask = (1 << shift) - 1

var spf [maxN + 1]int
var primeCnt [maxN + 1]int
var squareCnt [maxN + 1]int
var primeUsed []int
var squareUsed []int

func initSPF() {
    for i := 2; i <= maxN; i++ {
        if spf[i] == 0 {
            spf[i] = i
            if i*i <= maxN {
                for j := i * i; j <= maxN; j += i {
                    if spf[j] == 0 {
                        spf[j] = i
                    }
                }
            }
        }
    }
    spf[0], spf[1] = 1, 1
}

type kindInfo struct {
    kind int
    p    int
    q    int
}

func categorize(x int) kindInfo {
    temp := x
    primes := [2]int{}
    exps := [2]int{}
    cnt := 0
    for temp > 1 {
        p := spf[temp]
        e := 0
        for temp%p == 0 {
            temp /= p
            e++
        }
        if cnt < 2 {
            primes[cnt] = p
            exps[cnt] = e
        }
        cnt++
        if cnt > 2 {
            return kindInfo{}
        }
    }
    if cnt == 0 {
        return kindInfo{}
    }
    if cnt == 1 {
        if exps[0] == 1 {
            return kindInfo{kind: 1, p: primes[0]}
        }
        if exps[0] == 2 {
            return kindInfo{kind: 2, p: primes[0]}
        }
        return kindInfo{}
    }
    if cnt == 2 {
        if exps[0] == 1 && exps[1] == 1 {
            p := primes[0]
            q := primes[1]
            if p > q {
                p, q = q, p
            }
            if p == q {
                return kindInfo{}
            }
            return kindInfo{kind: 3, p: p, q: q}
        }
    }
    return kindInfo{}
}

func main() {
    initSPF()
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        semiMap := make(map[int64]int)
        for i := 0; i < n; i++ {
            var x int
            fmt.Fscan(in, &x)
            info := categorize(x)
            switch info.kind {
            case 1:
                if primeCnt[info.p] == 0 {
                    primeUsed = append(primeUsed, info.p)
                }
                primeCnt[info.p]++
            case 2:
                if squareCnt[info.p] == 0 {
                    squareUsed = append(squareUsed, info.p)
                }
                squareCnt[info.p]++
            case 3:
                key := (int64(info.p) << shift) | int64(info.q)
                semiMap[key]++
            default:
                // ignore numbers that can't be part of valid pair
            }
        }

        var totalPrime int64
        var sumSquares int64
        for _, p := range primeUsed {
            c := int64(primeCnt[p])
            totalPrime += c
            sumSquares += c * c
        }
        ans := (totalPrime*totalPrime - sumSquares) / 2

        for key, cnt := range semiMap {
            p := int(key >> shift)
            q := int(key & mask)
            c := int64(cnt)
            a := int64(primeCnt[p])
            b := int64(primeCnt[q])
            ans += a*c + b*c + c*(c+1)/2
        }

        for _, p := range squareUsed {
            b := int64(squareCnt[p])
            a := int64(primeCnt[p])
            ans += b*(b+1)/2 + a*b
        }

        fmt.Fprintln(out, ans)

        for _, p := range primeUsed {
            primeCnt[p] = 0
        }
        primeUsed = primeUsed[:0]
        for _, p := range squareUsed {
            squareCnt[p] = 0
        }
        squareUsed = squareUsed[:0]
    }
}
