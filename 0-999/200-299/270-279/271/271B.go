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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([][]int, n)
   maxA := 0
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &a[i][j])
           if a[i][j] > maxA {
               maxA = a[i][j]
           }
       }
   }
   // Sieve up to maxA + buffer
   buf := 500
   limit := maxA + buf
   isPrime := make([]bool, limit+1)
   for i := 2; i <= limit; i++ {
       isPrime[i] = true
   }
   for i := 2; i*i <= limit; i++ {
       if isPrime[i] {
           for j := i * i; j <= limit; j += i {
               isPrime[j] = false
           }
       }
   }
   // nextPrime[x] = smallest prime >= x
   nextPrime := make([]int, limit+2)
   next := -1
   for i := limit; i >= 0; i-- {
       if isPrime[i] {
           next = i
       }
       nextPrime[i] = next
   }
   // compute row and column sums
   rowSum := make([]int, n)
   colSum := make([]int, m)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           v := a[i][j]
           np := nextPrime[v]
           if np < 0 {
               // should not happen
               np = v
           }
           delta := np - v
           rowSum[i] += delta
           colSum[j] += delta
       }
   }
   // find minimum
   ans := rowSum[0]
   for i := 0; i < n; i++ {
       if rowSum[i] < ans {
           ans = rowSum[i]
       }
   }
   for j := 0; j < m; j++ {
       if colSum[j] < ans {
           ans = colSum[j]
       }
   }
   fmt.Fprintln(writer, ans)
}
