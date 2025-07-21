package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read all integers from input: supplies for ingredients
   // Order: carrots, calories, chocolate (g), flour (packs), egg
   var vals []int
   for {
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           break
       }
       vals = append(vals, x)
   }
   if len(vals) < 5 {
       // Not enough data, nothing to do
       return
   }
   carrots := vals[0]
   // calories := vals[1]  // not used (recipe uses 0 calories)
   chocolate := vals[2]
   flour := vals[3]
   egg := vals[4]

   // Recipe needs per cake:
   // 2 carrots, 0 calories, 100 g chocolate, 1 pack flour, 1 egg
   maxByCarrots := carrots / 2
   maxByChocolate := chocolate / 100
   maxByFlour := flour
   maxByEgg := egg

   // Compute minimum of these
   result := maxByCarrots
   if maxByChocolate < result {
       result = maxByChocolate
   }
   if maxByFlour < result {
       result = maxByFlour
   }
   if maxByEgg < result {
       result = maxByEgg
   }
   fmt.Println(result)
}
