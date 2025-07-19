package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var Q int
   fmt.Fscan(reader, &Q)
   for q := 0; q < Q; q++ {
       var N int
       fmt.Fscan(reader, &N)
       a := make([]int, N)
       for i := 0; i < N; i++ {
           fmt.Fscan(reader, &a[i])
       }
       ans := solve(N, a)
       fmt.Fprintln(writer, len(ans))
       for _, v := range ans {
           fmt.Fprint(writer, len(v))
           for _, val := range v {
               fmt.Fprint(writer, " ", val)
           }
           fmt.Fprintln(writer)
       }
   }
}

// solve decomposes permutation a of length N into sequences
func solve(N int, a []int) [][]int {
   covered := make([]bool, N)
   dp := make([]int, N)
   x := make([]int, N+1)
   pre := make([]int, N+1)
   valToPos := make([]int, N+1)
   for i := 0; i < N; i++ {
       valToPos[a[i]] = i
       covered[i] = false
       dp[i] = i + 1
   }
   // initial D: minimal D such that (D+1)*(D+2)/2 - 1 >= N
   D := 1
   for (D+1)*(D+2)/2-1 < N {
       D++
   }
   ans := [][]int{}
   rem := N
   preCut := N + 5
   for {
       tmp := main3(N, a, covered, dp, x, pre, valToPos, D, preCut, &ans)
       if tmp == -1 {
           break
       }
       D--
       rem -= tmp
       if rem == 0 {
           break
       }
       preCut = tmp
   }
   return ans
}

// main3 performs one step of extracting a sequence or final grouping
func main3(N int, a []int, covered []bool, dp []int, x []int, pre []int, valToPos []int, D int, preCut int, ans *[][]int) int {
   INF := N + 5
   for i := 0; i <= N; i++ {
       x[i] = INF
   }
   x[0] = 0
   for i := 0; i < N; i++ {
       if !covered[i] {
           high := dp[i]
           low := max(dp[i]-preCut, 1) - 1
           for high-low > 1 {
               mid := (high + low) / 2
               if x[mid] < a[i] {
                   low = mid
               } else {
                   high = mid
               }
           }
           dp[i] = high
           pre[a[i]] = x[low]
           x[high] = a[i]
       }
   }
   maxdp, maxPos := 0, -1
   for i := 0; i < N; i++ {
       if !covered[i] && dp[i] > maxdp {
           maxdp = dp[i]
           maxPos = i
       }
   }
   if maxdp <= D {
       base := len(*ans)
       for i := 0; i < maxdp; i++ {
           *ans = append(*ans, []int{})
       }
       for i := 0; i < N; i++ {
           if !covered[i] {
               idx := base + dp[i] - 1
               (*ans)[idx] = append((*ans)[idx], a[i])
           }
       }
       return -1
   }
   y := a[maxPos]
   seq := []int{}
   for y > 0 {
       seq = append(seq, y)
       covered[valToPos[y]] = true
       y = pre[y]
   }
   // reverse seq
   for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
       seq[i], seq[j] = seq[j], seq[i]
   }
   *ans = append(*ans, seq)
   return len(seq)
}

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
