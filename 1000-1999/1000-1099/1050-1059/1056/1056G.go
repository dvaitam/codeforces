package main

import (
   "bufio"
   "fmt"
   "os"
)

func solve(nn, mm, ss int, tt int64) int {
   n, m := nn, mm
   s := ss - 1
   t := tt
   been := make([]int64, n+5)
   // initial moves until t % n == 0
   for t % int64(n) != 0 {
       move := int(t % int64(n))
       if s < m {
           s += move
       } else {
           s -= move
       }
       // wrap around
       s = ((s % n) + n) % n
       t--
   }
   // number of full cycles
   times := t / int64(n)
   for times > 0 {
       if been[s] != 0 {
           cycleLen := been[s] - times
           if cycleLen != 0 {
               times %= cycleLen
           }
           break
       }
       been[s] = times
       // perform one full cycle
       for i := n - 1; i > 0; i-- {
           if s < m {
               s += i
           } else {
               s -= i
           }
           // wrap around
           if s < 0 {
               s += n
           } else if s >= n {
               s -= n
           }
       }
       times--
   }
   // if remaining cycles lead to a specific earlier state
   if times != 0 {
       for i := 0; i < n; i++ {
           if been[i] != 0 && been[s] - been[i] == times {
               s = i
               break
           }
       }
   }
   return s + 1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var nn, mm, ss int
   var tt int64
   fmt.Fscan(reader, &nn, &mm, &ss, &tt)
   res := solve(nn, mm, ss, tt)
   fmt.Fprintln(writer, res)
}
