package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAXN = 100

var (
   dp [MAXN+1][MAXN+1][MAXN+1]int64
   n  int
   s  string
   A  []int64
   B  [MAXN]int
   m  int
   reader *bufio.Reader
   writer *bufio.Writer
)

// rec computes maximum profit for segments [i, j) with k extra characters
func rec(i, j, k int) int64 {
   if dp[i][j][k] != -1 {
       return dp[i][j][k]
   }
   // no segments: profit from k-length
   if j-i == 0 {
       return A[k]
   }
   // single segment: combine k with B[i]
   if j-i == 1 {
       return A[B[i]+k]
   }
   // option: take last segment alone
   ans := rec(i, j-1, 0) + A[B[j-1]+k]
   // option: merge last segment with a previous one, skipping even-length blocks
   for l := j-2; l >= i; l -= 2 {
       v := rec(i, l, B[j-1]+k) + rec(l, j-1, 0)
       if v > ans {
           ans = v
       }
   }
   dp[i][j][k] = ans
   return ans
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &s)
   A = make([]int64, n+1)
   // read base profits and compute optimal for each length
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &A[i])
       for j := 1; j < i; j++ {
           if A[j]+A[i-j] > A[i] {
               A[i] = A[j] + A[i-j]
           }
       }
   }
   // initialize dp array to -1
   for i := 0; i <= n; i++ {
       for j := 0; j <= n; j++ {
           for k := 0; k <= n; k++ {
               dp[i][j][k] = -1
           }
       }
   }
   // build segments of consecutive equal characters
   m = 0
   for i := 0; i < n; {
       j := i
       for j < n && s[i] == s[j] {
           j++
       }
       B[m] = j - i
       m++
       i = j
   }
   // compute and print result
   res := rec(0, m, 0)
   fmt.Fprintln(writer, res)
}
