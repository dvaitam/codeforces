package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var supply, demand string
   if _, err := fmt.Fscan(reader, &supply); err != nil {
       return
   }
   // Read desired garland pieces colors
   if _, err := fmt.Fscan(reader, &demand); err != nil {
       return
   }
   // Count available sheets per color
   supplyCount := make([]int, 26)
   for _, c := range supply {
       supplyCount[c-'a']++
   }
   // Track which colors are needed
   needed := make([]bool, 26)
   for _, c := range demand {
       needed[c-'a'] = true
   }
   // Compute total area
   total := 0
   for i, use := range needed {
       if use {
           if supplyCount[i] == 0 {
               fmt.Println(-1)
               return
           }
           total += supplyCount[i]
       }
   }
   fmt.Println(total)
}
