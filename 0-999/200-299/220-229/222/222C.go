package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAX = 10000001

var lp []int32
var primes []int
var factors = make(map[int]int)
var factors2 = make(map[int]int)

func factorSieve() {
   lp = make([]int32, MAX)
   lp[1] = 1
   for i := 2; i < MAX; i++ {
       if lp[i] == 0 {
           lp[i] = int32(i)
           primes = append(primes, i)
       }
       for j := 0; j < len(primes) && primes[j] <= int(lp[i]); j++ {
           x := i * primes[j]
           if x >= MAX {
               break
           }
           lp[x] = int32(primes[j])
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   factorSieve()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   b := make([]int, m)
   // Read a and count prime factors
   for i := 0; i < n; i++ {
       var x int
       fmt.Fscan(reader, &x)
       a[i] = x
       for x != 1 {
           y := int(lp[x])
           for x%y == 0 {
               x /= y
               factors[y]++
           }
       }
   }
   // Read b and count prime factors
   for i := 0; i < m; i++ {
       var x int
       fmt.Fscan(reader, &x)
       b[i] = x
       for x != 1 {
           y := int(lp[x])
           for x%y == 0 {
               x /= y
               factors2[y]++
           }
       }
   }
   // Output sizes
   fmt.Fprintf(writer, "%d %d\n", n, m)
   // Process a: remove common factors
   for i := 0; i < n; i++ {
       x := a[i]
       res := x
       for res != 1 {
           y := int(lp[res])
           for res%y == 0 {
               res /= y
               if factors2[y] > 0 {
                   x /= y
                   factors2[y]--
               }
           }
       }
       fmt.Fprintf(writer, "%d ", x)
   }
   fmt.Fprint(writer, "\n")
   // Process b: remove common factors
   for i := 0; i < m; i++ {
       res := b[i]
       for res != 1 {
           y := int(lp[res])
           for res%y == 0 {
               res /= y
               if factors[y] > 0 {
                   b[i] /= y
                   factors[y]--
               }
           }
       }
       fmt.Fprintf(writer, "%d ", b[i])
   }
   fmt.Fprint(writer, "\n")
}
