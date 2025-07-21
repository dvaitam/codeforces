package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF64 = int64(4e18)

var (
   n, k       int
   u          [][]int
   dpPrev, dpCur []int64
   curL, curR int
   curCost    int64
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // fast integer reading
   n = readInt(reader)
   k = readInt(reader)
   u = make([][]int, n+2)
   for i := 0; i <= n+1; i++ {
       u[i] = make([]int, n+2)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           u[i][j] = readInt(reader)
       }
   }
   dpPrev = make([]int64, n+1)
   dpCur = make([]int64, n+1)
   for i := 1; i <= n; i++ {
       dpPrev[i] = INF64
       dpCur[i] = INF64
   }
   dpPrev[0] = 0
   curL, curR, curCost = 1, 0, 0
   // DP layers
   for i := 1; i <= k; i++ {
       // reset dpCur for this layer
       for j := 0; j <= n; j++ {
           dpCur[j] = INF64
       }
       solve(i, n, 0, n-1)
       // swap dpPrev and dpCur
       dpPrev, dpCur = dpCur, dpPrev
   }
   // answer is dpPrev[n]
   fmt.Println(dpPrev[n])
}

func solve(l, r, optL, optR int) {
   if l > r {
       return
   }
   mid := (l + r) >> 1
   bestK := -1
   // search optimal k
   maxK := optR
   if mid-1 < maxK {
       maxK = mid - 1
   }
   for kx := optL; kx <= maxK; kx++ {
       c := dpPrev[kx] + costLR(kx+1, mid)
       if c < dpCur[mid] {
           dpCur[mid] = c
           bestK = kx
       }
   }
   if l == r {
       return
   }
   // divide and conquer
   if bestK < 0 {
       bestK = optL
   }
   solve(l, mid-1, optL, bestK)
   solve(mid+1, r, bestK, optR)
}

func costLR(L, R int) int64 {
   // move curL, curR to L, R
   for curL > L {
       curL--
       for i := curL + 1; i <= curR; i++ {
           curCost += int64(u[curL][i])
       }
   }
   for curR < R {
       curR++
       for i := curL; i < curR; i++ {
           curCost += int64(u[i][curR])
       }
   }
   for curL < L {
       for i := curL + 1; i <= curR; i++ {
           curCost -= int64(u[curL][i])
       }
       curL++
   }
   for curR > R {
       for i := curL; i < curR; i++ {
           curCost -= int64(u[i][curR])
       }
       curR--
   }
   return curCost
}

func readInt(r *bufio.Reader) int {
   var x int
   var c byte
   var neg bool
   // skip non-digit
   for {
       b, err := r.ReadByte()
       if err != nil {
           return x
       }
       c = b
       if c == '-' {
           neg = true
           break
       }
       if c >= '0' && c <= '9' {
           x = int(c - '0')
           break
       }
   }
   // read digits
   for {
       b, err := r.ReadByte()
       if err != nil {
           break
       }
       c = b
       if c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
   }
   if neg {
       return -x
   }
   return x
}
