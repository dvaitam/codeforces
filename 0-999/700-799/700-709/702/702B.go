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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   arr := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
   }

   // Precompute powers of two up to 2^31
   pow2 := make([]int64, 32)
   for i := 0; i < 32; i++ {
       pow2[i] = 1 << i
   }

   countMap := make(map[int64]int)
   var result int64
   for _, a := range arr {
       for _, p := range pow2 {
           need := p - a
           if cnt, ok := countMap[need]; ok {
               result += int64(cnt)
           }
       }
       countMap[a]++
   }

   fmt.Fprintln(writer, result)
}
