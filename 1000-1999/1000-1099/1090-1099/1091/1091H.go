package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   MAX  = 200200
   MEXN = 1000
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Sieve and prime factor preparation
   isp := make([]bool, MAX)
   pr := make([]int, MAX)
   cnt := make([]int, MAX)
   for i := 2; i < MAX; i++ {
       pr[i] = i
   }
   var moves []int
   // mark composites up to sqrt(MAX) and collect prime squares
   for i := 2; i*i < MAX; i++ {
       if !isp[i] {
           for j := i * i; j < MAX; j += i {
               isp[j] = true
           }
           moves = append(moves, i*i)
       }
   }
   // collect primes and divide out distinct prime factors
   for i := 2; i < MAX; i++ {
       if !isp[i] {
           moves = append(moves, i)
           for j := i; j < MAX; j += i {
               pr[j] /= i
               cnt[j]++
           }
       }
   }
   // read input
   var N, F int
   fmt.Fscan(reader, &N, &F)
   // collect semiprimes except F
   for i := 2; i < MAX; i++ {
       if pr[i] == 1 && cnt[i] == 2 && i != F {
           moves = append(moves, i)
       }
   }
   sort.Ints(moves)
   chk := make([]bool, MAX)
   for _, x := range moves {
       if x < MAX {
           chk[x] = true
       }
   }

   // compute mex values
   mex := make([]int, MAX)
   in := make([][]int, MEXN)
   in[0] = append(in[0], 0)
   for i := 1; i < MAX; i++ {
       // find smallest mex not reachable
       for {
           g := false
           mi := mex[i]
           for _, p := range in[mi] {
               if i-p >= 0 && chk[i-p] {
                   g = true
                   break
               }
           }
           if !g {
               break
           }
           mex[i]++
       }
       mi := mex[i]
       if mi < MEXN {
           in[mi] = append(in[mi], i)
       }
   }

   // game result
   xor := 0
   for k := 0; k < N; k++ {
       var a, b, c int
       fmt.Fscan(reader, &a, &b, &c)
       xor ^= mex[b-a-1]
       xor ^= mex[c-b-1]
   }
   if xor != 0 {
       fmt.Fprintln(writer, "Alice")
       fmt.Fprintln(writer, "Bob")
   } else {
       fmt.Fprintln(writer, "Bob")
       fmt.Fprintln(writer, "Alice")
   }
}
