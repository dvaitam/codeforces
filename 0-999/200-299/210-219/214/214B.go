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
   freq := make([]int, 10)
   for i := 0; i < n; i++ {
       var d int
       fmt.Fscan(reader, &d)
       if d >= 0 && d <= 9 {
           freq[d]++
       }
   }
   // must have at least one zero for divisibility by 10
   if freq[0] == 0 {
       fmt.Fprint(writer, -1)
       return
   }
   // sum of digits
   sum := 0
   for d, c := range freq {
       sum += d * c
   }
   mod := sum % 3
   // prepare lists of digits by mod class
   mod1 := make([]int, 0, n)
   mod2 := make([]int, 0, n)
   for d := 1; d <= 9; d++ {
       for i := 0; i < freq[d]; i++ {
           if d%3 == 1 {
               mod1 = append(mod1, d)
           } else if d%3 == 2 {
               mod2 = append(mod2, d)
           }
       }
   }
   // helper to remove digits
   remove := func(list []int, cnt int) bool {
       if len(list) < cnt {
           return false
       }
       // remove smallest cnt digits
       for i := 0; i < cnt; i++ {
           digit := list[i]
           freq[digit]--
       }
       return true
   }
   switch mod {
   case 1:
       if !remove(mod1, 1) {
           if !remove(mod2, 2) {
               fmt.Fprint(writer, -1)
               return
           }
       }
   case 2:
       if !remove(mod2, 1) {
           if !remove(mod1, 2) {
               fmt.Fprint(writer, -1)
               return
           }
       }
   }
   // after removal, ensure zero still exists
   if freq[0] == 0 {
       fmt.Fprint(writer, -1)
       return
   }
   // check if only zeros remain
   nonZero := 0
   for d := 1; d <= 9; d++ {
       nonZero += freq[d]
   }
   if nonZero == 0 {
       // only zeros
       fmt.Fprint(writer, 0)
       return
   }
   // build result: digits in descending order
   for d := 9; d >= 0; d-- {
       for i := 0; i < freq[d]; i++ {
           writer.WriteByte(byte('0' + d))
       }
   }
}
