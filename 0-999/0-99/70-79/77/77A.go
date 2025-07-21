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
   // map hero names to indices 0..6
   names := []string{"Anka", "Chapay", "Cleo", "Troll", "Dracul", "Snowy", "Hexadecimal"}
   nameToIdx := make(map[string]int, len(names))
   for i, s := range names {
       nameToIdx[s] = i
   }
   // read likes
   var edges [][2]int
   for i := 0; i < n; i++ {
       var p, likesWord, q string
       fmt.Fscan(reader, &p, &likesWord, &q)
       pi := nameToIdx[p]
       qi := nameToIdx[q]
       edges = append(edges, [2]int{pi, qi})
   }
   // read experience values
   var a, b, c int64
   fmt.Fscan(reader, &a, &b, &c)
   xpValues := [3]int64{a, b, c}

   bestDiff := int64(1<<62 - 1)
   bestLikes := 0
   // assignment: for each hero, team 0,1,2
   var assign [7]int
   for mask := 0; mask < 2187; mask++ {
       m := mask
       var sz [3]int
       for i := 0; i < 7; i++ {
           assign[i] = m % 3
           sz[assign[i]]++
           m /= 3
       }
       // each team must have at least one member
       if sz[0] == 0 || sz[1] == 0 || sz[2] == 0 {
           continue
       }
       // compute xp per group
       var xpGroup [3]int64
       for i := 0; i < 3; i++ {
           xpGroup[i] = xpValues[i] / int64(sz[i])
       }
       // compute min and max xp
       var minXP, maxXP int64
       minXP = xpGroup[assign[0]]
       maxXP = minXP
       for i := 1; i < 7; i++ {
           x := xpGroup[assign[i]]
           if x < minXP {
               minXP = x
           }
           if x > maxXP {
               maxXP = x
           }
       }
       diff := maxXP - minXP
       if diff > bestDiff {
           continue
       }
       // compute likes for this partition
       cntLikes := 0
       for _, e := range edges {
           if assign[e[0]] == assign[e[1]] {
               cntLikes++
           }
       }
       if diff < bestDiff {
           bestDiff = diff
           bestLikes = cntLikes
       } else if cntLikes > bestLikes {
           bestLikes = cntLikes
       }
   }
   fmt.Println(bestDiff, bestLikes)
}
