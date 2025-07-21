package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // Track unique letters
   used := make([]bool, 24)
   masks := make([]uint32, n)
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       var m uint32
       for _, ch := range s {
           idx := ch - 'a'
           if idx < 0 || idx >= 24 {
               continue
           }
           used[idx] = true
           m |= 1 << idx
       }
       masks[i] = m
   }
   // If not all 24 letters used, each value repeats even times => XOR = 0
   cntUsed := 0
   for _, u := range used {
       if u {
           cntUsed++
       }
   }
   if cntUsed < 24 {
       fmt.Println(0)
       return
   }
   // f[mask] = count of words with exactly that unique-letter mask
   size := 1 << 24
   f := make([]uint16, size)
   for _, m := range masks {
       f[m]++
   }
   // SOS DP: f[mask] = sum_{submask} f[submask]
   for bit := 0; bit < 24; bit++ {
       step := 1 << bit
       for mask := 0; mask < size; mask++ {
           if mask&(step) != 0 {
               f[mask] += f[mask^step]
           }
       }
   }
   // Compute XOR of squares
   var res uint32
   for mask := 0; mask < size; mask++ {
       v := uint32(f[mask])
       res ^= v * v
   }
   fmt.Println(res)
}
