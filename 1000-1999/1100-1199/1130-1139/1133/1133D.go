package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(x int64) int64 {
   if x < 0 {
       return -x
   }
   return x
}

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   b := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }

   alwaysZero := 0
   counts := make(map[[2]int64]int)
   for i := 0; i < n; i++ {
       ai, bi := a[i], b[i]
       if ai == 0 {
           if bi == 0 {
               alwaysZero++
           }
           continue
       }
       num := -bi
       den := ai
       if den < 0 {
           den = -den
           num = -num
       }
       if num == 0 {
           den = 1
       } else {
           g := gcd(abs(num), den)
           num /= g
           den /= g
       }
       key := [2]int64{num, den}
       counts[key]++
   }

   maxCount := 0
   for _, v := range counts {
       if v > maxCount {
           maxCount = v
       }
   }
   result := alwaysZero + maxCount
   fmt.Fprintln(writer, result)
}
