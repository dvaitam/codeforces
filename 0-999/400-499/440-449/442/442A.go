package main

import (
   "fmt"
   "math/bits"
)

func colorIndex(c byte) int {
   switch c {
   case 'R': return 0
   case 'G': return 1
   case 'B': return 2
   case 'Y': return 3
   case 'W': return 4
   }
   return -1
}

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   // Record present card types
   present := make([][]bool, 5)
   for i := 0; i < 5; i++ {
       present[i] = make([]bool, 5)
   }
   var s string
   for i := 0; i < n; i++ {
       fmt.Scan(&s)
       ci := colorIndex(s[0])
       vi := int(s[1] - '1')
       if ci >= 0 && vi >= 0 && ci < 5 && vi < 5 {
           present[ci][vi] = true
       }
   }
   // Build slice of types
   types := make([][2]int, 0, 25)
   for ci := 0; ci < 5; ci++ {
       for vi := 0; vi < 5; vi++ {
           if present[ci][vi] {
               types = append(types, [2]int{ci, vi})
           }
       }
   }
   best := 10
   // Try all subsets of hints (5 colors + 5 values)
   for mask := 0; mask < (1 << 10); mask++ {
       cnt := bits.OnesCount(uint(mask))
       if cnt >= best {
           continue
       }
       ok := true
       // Check all pairs of distinct types
       for i := 0; i < len(types) && ok; i++ {
           for j := i + 1; j < len(types); j++ {
               c1, v1 := types[i][0], types[i][1]
               c2, v2 := types[j][0], types[j][1]
               distinguished := false
               for bit := 0; bit < 10; bit++ {
                   if mask&(1<<bit) == 0 {
                       continue
                   }
                   if bit < 5 {
                       // color hint
                       if (c1 == bit) != (c2 == bit) {
                           distinguished = true
                           break
                       }
                   } else {
                       // value hint
                       vb := bit - 5
                       if (v1 == vb) != (v2 == vb) {
                           distinguished = true
                           break
                       }
                   }
               }
               if !distinguished {
                   ok = false
                   break
               }
           }
       }
       if ok {
           best = cnt
       }
   }
   fmt.Println(best)
}
