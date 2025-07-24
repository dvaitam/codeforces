package main

import (
   "bufio"
   "fmt"
   "os"
)

type pair struct { r, b int }

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   isRed := make([]bool, n)
   rreq := make([]int, n)
   breq := make([]int, n)
   for i := 0; i < n; i++ {
       var c byte
       fmt.Fscan(in, &c, &rreq[i], &breq[i])
       isRed[i] = (c == 'R')
   }
   // precompute counts for masks
   N := 1 << n
   rcnt := make([]int, N)
   bcnt := make([]int, N)
   for mask := 1; mask < N; mask++ {
       lsb := mask & -mask
       i := bitsTrailing(lsb)
       prev := mask ^ lsb
       if isRed[i] {
           rcnt[mask] = rcnt[prev] + 1
           bcnt[mask] = bcnt[prev]
       } else {
           bcnt[mask] = bcnt[prev] + 1
           rcnt[mask] = rcnt[prev]
       }
   }
   // sum of reqs as upper bound
   sumR, sumB := 0, 0
   for i := 0; i < n; i++ {
       sumR += rreq[i]
       sumB += breq[i]
   }
   lo, hi := 0, sumR+sumB
   goal := N - 1
   for lo < hi {
       mid := (lo + hi) / 2
       if can(mid, n, rreq, breq, isRed, rcnt, bcnt, goal) {
           hi = mid
       } else {
           lo = mid + 1
       }
   }
   // total turns = collects + n buys
   fmt.Println(lo + n)
}

// bitsTrailing returns index of the single bit in v (v is power of two)
func bitsTrailing(v int) int {
   // simple loop since n<=16
   for i := 0; i < 32; i++ {
       if v>>i&1 == 1 {
           return i
       }
   }
   return 0
}

// can check if with t collects we can buy all cards
func can(t, n int, rreq, breq []int, isRed []bool, rcnt, bcnt []int, goal int) bool {
   N := 1 << n
   dp := make([][]pair, N)
   dp[0] = []pair{{t, t}}
   for mask := 0; mask < N; mask++ {
       states := dp[mask]
       if len(states) == 0 {
           continue
       }
       // prune dominated states
       dp[mask] = prune(states)
       if mask == goal && len(dp[mask]) > 0 {
           return true
       }
       A := rcnt[mask]
       B := bcnt[mask]
       for _, st := range dp[mask] {
           rrem, brem := st.r, st.b
           for i := 0; i < n; i++ {
               bit := 1 << i
               if mask&bit != 0 {
                   continue
               }
               needR := rreq[i] - A
               if needR < 0 {
                   needR = 0
               }
               needB := breq[i] - B
               if needB < 0 {
                   needB = 0
               }
               if rrem >= needR && brem >= needB {
                   nm := mask | bit
                   dp[nm] = append(dp[nm], pair{rrem - needR, brem - needB})
               }
           }
       }
   }
   return len(dp[goal]) > 0
}

// prune removes dominated states: keeps sorted by r desc, and only those with strictly larger b
func prune(arr []pair) []pair {
   // sort by r desc, b desc
   // simple insertion sort since small
   for i := 1; i < len(arr); i++ {
       j := i
       for j > 0 && (arr[j].r > arr[j-1].r || (arr[j].r == arr[j-1].r && arr[j].b > arr[j-1].b)) {
           arr[j], arr[j-1] = arr[j-1], arr[j]
           j--
       }
   }
   res := make([]pair, 0, len(arr))
   maxB := -1
   for _, p := range arr {
       if p.b > maxB {
           res = append(res, p)
           maxB = p.b
       }
   }
   return res
}
