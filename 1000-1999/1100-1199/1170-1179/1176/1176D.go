package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   total := 2 * n
   vals := make([]int, total)
   maxVal := 0
   for i := 0; i < total; i++ {
       fmt.Fscan(reader, &vals[i])
       if vals[i] > maxVal {
           maxVal = vals[i]
       }
   }
   // Sieve for primes, spf, and prime positions
   N := maxVal
   spf := make([]int, N+1)
   primes := make([]int, 0, N/10)
   primePos := make([]int, N+1)
   for i := 2; i <= N; i++ {
       if spf[i] == 0 {
           spf[i] = i
           primePos[i] = len(primes) + 1
           primes = append(primes, i)
       }
       for _, p := range primes {
           if p > spf[i] || i*p > N {
               break
           }
           spf[i*p] = p
       }
   }
   // Partition values
   a := make([]int, 0, total)
   b := make([]int, 0, total)
   isPrime := make([]bool, N+1)
   for _, p := range primes {
       isPrime[p] = true
   }
   for _, v := range vals {
       if v <= N && isPrime[v] {
           a = append(a, v)
       } else {
           b = append(b, v)
       }
   }
   // Prepare low counts and answer
   low := make([]int, N+1)
   ans := make([]int, 0, n)

   // Process composites b in descending order
   sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
   for _, v := range b {
       if low[v] > 0 {
           low[v]--
       } else {
           ans = append(ans, v)
           u := v / spf[v]
           low[u]++
       }
   }
   // Process primes a in ascending order
   sort.Ints(a)
   for _, v := range a {
       if low[v] > 0 {
           low[v]--
       } else {
           // use prime position
           ans = append(ans, primePos[v])
           low[primePos[v]]++
       }
   }
   // Output first n answers
   for i := 0; i < n; i++ {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, ans[i])
   }
   writer.WriteByte('\n')
}
