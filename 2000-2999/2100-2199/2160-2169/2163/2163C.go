package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

type fastScanner struct {
    r *bufio.Reader
}

func newFastScanner() *fastScanner {
    return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
    sign, val := 1, 0
    c, _ := fs.r.ReadByte()
    for (c < '0' || c > '9') && c != '-' {
        c, _ = fs.r.ReadByte()
    }
    if c == '-' {
        sign = -1
        c, _ = fs.r.ReadByte()
    }
    for c >= '0' && c <= '9' {
        val = val*10 + int(c-'0')
        c, _ = fs.r.ReadByte()
    }
    return sign * val
}

func main() {
    fs := newFastScanner()
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    t := fs.nextInt()
    for ; t > 0; t-- {
        n := fs.nextInt()
        top := make([]int, n)
        bottom := make([]int, n)
        for i := 0; i < n; i++ {
            top[i] = fs.nextInt()
        }
        for i := 0; i < n; i++ {
            bottom[i] = fs.nextInt()
        }

        prefMin := make([]int, n)
        prefMax := make([]int, n)
        prefMin[0], prefMax[0] = top[0], top[0]
        for i := 1; i < n; i++ {
            if top[i] < prefMin[i-1] {
                prefMin[i] = top[i]
            } else {
                prefMin[i] = prefMin[i-1]
            }
            if top[i] > prefMax[i-1] {
                prefMax[i] = top[i]
            } else {
                prefMax[i] = prefMax[i-1]
            }
        }

        sufMin := make([]int, n)
        sufMax := make([]int, n)
        sufMin[n-1], sufMax[n-1] = bottom[n-1], bottom[n-1]
        for i := n - 2; i >= 0; i-- {
            if bottom[i] < sufMin[i+1] {
                sufMin[i] = bottom[i]
            } else {
                sufMin[i] = sufMin[i+1]
            }
            if bottom[i] > sufMax[i+1] {
                sufMax[i] = bottom[i]
            } else {
                sufMax[i] = sufMax[i+1]
            }
        }

        type pair struct{ R, L int }
        pairs := make([]pair, n)
        for i := 0; i < n; i++ {
            L := prefMin[i]
            if sufMin[i] < L {
                L = sufMin[i]
            }
            R := prefMax[i]
            if sufMax[i] > R {
                R = sufMax[i]
            }
            pairs[i] = pair{R: R, L: L}
        }

        sort.Slice(pairs, func(i, j int) bool {
            if pairs[i].R == pairs[j].R {
                return pairs[i].L > pairs[j].L
            }
            return pairs[i].R < pairs[j].R
        })

        ans := int64(0)
        curMaxL, idx := 0, 0
        total := 2 * n
        for r := 1; r <= total; r++ {
            for idx < n && pairs[idx].R <= r {
                if pairs[idx].L > curMaxL {
                    curMaxL = pairs[idx].L
                }
                idx++
            }
            ans += int64(curMaxL)
        }

        fmt.Fprintln(out, ans)
    }
}
