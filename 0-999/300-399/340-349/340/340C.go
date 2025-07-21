package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
   var sumA int64
   for _, v := range a {
       sumA += v
   }
   var prefix int64
   var sumPairs int64
   for i, v := range a {
       // sum of |v - previous elements|
       sumPairs += v*int64(i) - prefix
       prefix += v
   }
   // expected distance = (sumA + 2*sumPairs) / n
   num := sumA + 2*sumPairs
   den := int64(n)
   g := gcd(num, den)
   num /= g
   den /= g
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%d %d", num, den)
}
