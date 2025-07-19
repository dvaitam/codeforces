package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var A int64
   fmt.Fscan(reader, &n, &A)

   dd := make([]int64, n+1)
   profit := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       var di, ci int64
       fmt.Fscan(reader, &di, &ci)
       dd[i] = di
       profit[i] = A - ci
   }

   // compute differences between consecutive difficulties
   d := make([]int64, n)
   for i := 1; i < n; i++ {
       d[i] = dd[i+1] - dd[i]
   }

   // initial answer: best single problem
   ans := profit[1]
   for i := 2; i <= n; i++ {
       if profit[i] > ans {
           ans = profit[i]
       }
   }

   // find boundaries where each gap is maximum
   a1 := make([]int, n)
   a2 := make([]int, n)
   const inf = int64(1e18)
   // left boundaries
   stackVals := []int64{inf}
   stackIdx := []int{0}
   for i := 1; i < n; i++ {
       for len(stackVals) > 0 && stackVals[len(stackVals)-1] <= d[i] {
           stackVals = stackVals[:len(stackVals)-1]
           stackIdx = stackIdx[:len(stackIdx)-1]
       }
       a1[i] = stackIdx[len(stackIdx)-1] + 1
       stackVals = append(stackVals, d[i])
       stackIdx = append(stackIdx, i)
   }
   // right boundaries
   stackVals = []int64{inf}
   stackIdx = []int{n}
   for i := n - 1; i >= 1; i-- {
       for len(stackVals) > 0 && stackVals[len(stackVals)-1] <= d[i] {
           stackVals = stackVals[:len(stackVals)-1]
           stackIdx = stackIdx[:len(stackIdx)-1]
       }
       a2[i] = stackIdx[len(stackIdx)-1]
       stackVals = append(stackVals, d[i])
       stackIdx = append(stackIdx, i)
   }

   // prefix sums of profits
   prefix := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       prefix[i] = prefix[i-1] + profit[i]
   }

   // build sparse tables for min and max prefix sums
   size := n + 1
   maxLog := bits.Len(uint(size))
   stMin := make([][]int64, maxLog)
   stMax := make([][]int64, maxLog)
   for k := 0; k < maxLog; k++ {
       stMin[k] = make([]int64, size)
       stMax[k] = make([]int64, size)
   }
   for i := 0; i < size; i++ {
       stMin[0][i] = prefix[i]
       stMax[0][i] = prefix[i]
   }
   for k := 1; k < maxLog; k++ {
       span := 1 << (k - 1)
       for i := 0; i+(1<<k) <= size; i++ {
           stMin[k][i] = min(stMin[k-1][i], stMin[k-1][i+span])
           stMax[k][i] = max(stMax[k-1][i], stMax[k-1][i+span])
       }
   }

   // helper functions for RMQ
   rmqMin := func(l, r int) int64 {
       length := r - l + 1
       k := bits.Len(uint(length)) - 1
       return min(stMin[k][l], stMin[k][r-(1<<k)+1])
   }
   rmqMax := func(l, r int) int64 {
       length := r - l + 1
       k := bits.Len(uint(length)) - 1
       return max(stMax[k][l], stMax[k][r-(1<<k)+1])
   }

   // compute best segment including gap costs
   for i := 1; i < n; i++ {
       maxP := rmqMax(i+1, a2[i])
       minP := rmqMin(a1[i]-1, i-1)
       val := maxP - minP - d[i]*d[i]
       if val > ans {
           ans = val
       }
   }

   fmt.Fprintln(writer, ans)
}
