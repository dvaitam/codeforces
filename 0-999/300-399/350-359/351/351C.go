package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var m int64
   fmt.Fscan(in, &n, &m)
   a := make([]int64, n)
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &b[i])
   }
   // Build transition cost matrix C of size (n+1)
   N := n + 1
   INF := int64(9e18)
   // initialize C[u][v]
   C := make([][]int64, N)
   for i := 0; i < N; i++ {
       C[i] = make([]int64, N)
       for j := 0; j < N; j++ {
           C[i][j] = INF
       }
   }
   // DP1 for u = 0..n-1
   for u := 0; u < n; u++ {
       maxh := n - u // cap relative h at this for overflow
       dp := make([]int64, maxh+1)
       ndp := make([]int64, maxh+1)
       for h := 0; h <= maxh; h++ {
           dp[h] = INF
       }
       dp[0] = 0
       for i := 0; i < n; i++ {
           for h := 0; h <= maxh; h++ {
               ndp[h] = INF
           }
           for h := 0; h <= maxh; h++ {
               cost := dp[h]
               if cost >= INF {
                   continue
               }
               // open
               h2 := h + 1
               if h2 > maxh {
                   h2 = maxh
               }
               ndp[h2] = min(ndp[h2], cost + a[i])
               // close
               if h > 0 {
                   h3 := h - 1
                   ndp[h3] = min(ndp[h3], cost + b[i])
               }
           }
           dp, ndp = ndp, dp
       }
       // map dp final to C[u][v]
       for h := 0; h <= maxh; h++ {
           cost := dp[h]
           if cost >= INF {
               continue
           }
           f := u + h
           v := f
           if f >= n {
               v = n
           }
           if C[u][v] > cost {
               C[u][v] = cost
           }
       }
   }
   // DP2 for u = n (>=n state)
   // h range [-n..n], cap both ends
   base := n
   sizeH := 2*n + 1
   off := n
   dp2 := make([]int64, sizeH)
   ndp2 := make([]int64, sizeH)
   for i := 0; i < sizeH; i++ {
       dp2[i] = INF
   }
   dp2[off] = 0 // h=0
   for i := 0; i < n; i++ {
       for j := 0; j < sizeH; j++ {
           ndp2[j] = INF
       }
       for hi := 0; hi < sizeH; hi++ {
           cost := dp2[hi]
           if cost >= INF {
               continue
           }
           h := hi - off
           // open
           h2 := h + 1
           if h2 > n {
               h2 = n
           }
           hi2 := h2 + off
           ndp2[hi2] = min(ndp2[hi2], cost + a[i])
           // close
           h3 := h - 1
           if h3 < -n {
               h3 = -n
           }
           hi3 := h3 + off
           ndp2[hi3] = min(ndp2[hi3], cost + b[i])
       }
       dp2, ndp2 = ndp2, dp2
   }
   // map dp2 to C[n][v]
   for hi := 0; hi < sizeH; hi++ {
       cost := dp2[hi]
       if cost >= INF {
           continue
       }
       h := hi - off
       f := base + h
       v := f
       if f < n {
           v = f
       } else {
           v = n
       }
       if C[n][v] > cost {
           C[n][v] = cost
       }
   }
   // min-plus matrix exponentiation of C^m
   // initialize result as identity (0 on diag, INF else)
   res := make([][]int64, N)
   for i := 0; i < N; i++ {
       res[i] = make([]int64, N)
       for j := 0; j < N; j++ {
           if i == j {
               res[i][j] = 0
           } else {
               res[i][j] = INF
           }
       }
   }
   // base matrix
   baseM := C
   for m > 0 {
       if m & 1 == 1 {
           res = mult(res, baseM, INF)
       }
       baseM = mult(baseM, baseM, INF)
       m >>= 1
   }
   // answer is res[0][0]
   ans := res[0][0]
   fmt.Println(ans)
}

// mult returns A*B under min-plus, matrices size NxN
func mult(A, B [][]int64, INF int64) [][]int64 {
   n := len(A)
   C := make([][]int64, n)
   for i := 0; i < n; i++ {
       C[i] = make([]int64, n)
       for j := 0; j < n; j++ {
           C[i][j] = INF
           for k := 0; k < n; k++ {
               v := A[i][k] + B[k][j]
               if v < C[i][j] {
                   C[i][j] = v
               }
           }
       }
   }
   return C
}
