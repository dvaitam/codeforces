package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
   "math/bits"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   lesha := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &lesha[i])
   }
   // map words to index
   idxMap := make(map[string]int, n)
   for i, w := range lesha {
       idxMap[w] = i
   }
   var m int
   fmt.Fscan(reader, &m)
   // precompute inversion parameters
   fullMask := 1<<n - 1
   invMax := n*(n-1)/2
   invSize := invMax + 1
   // highMask[i] has bits > i
   highMask := make([]int, n)
   for i := 0; i < n; i++ {
       if i+1 < n {
           highMask[i] = ((fullMask) ^ ((1 << (i+1)) - 1))
       } else {
           highMask[i] = 0
       }
   }
   const INF = int32(1000000)
   totalStates := (fullMask + 1) * invSize
   dp := make([]int32, totalStates)
   bestIdx := -1
   bestP := -1
   for prob := 1; prob <= m; prob++ {
       var k int
       fmt.Fscan(reader, &k)
       // positions of each Lesha word in archive
       pos := make([][]int, n)
       for j := 0; j < k; j++ {
           var w string
           fmt.Fscan(reader, &w)
           if ii, ok := idxMap[w]; ok {
               pos[ii] = append(pos[ii], j)
           }
       }
       // init dp
       for i := range dp {
           dp[i] = INF
       }
       dp[0*invSize+0] = -1
       // dp over masks
       for mask := 0; mask <= fullMask; mask++ {
           base := mask * invSize
           for x := 0; x <= invMax; x++ {
               curPos := dp[base+x]
               if curPos >= INF {
                   continue
               }
               // try next word i
               for i := 0; i < n; i++ {
                   bit := 1 << i
                   if mask&bit != 0 {
                       continue
                   }
                   list := pos[i]
                   // find next occurrence
                   // curPos is int32, convert
                   cp := int(curPos)
                   j := sort.Search(len(list), func(j int) bool { return list[j] > cp })
                   if j >= len(list) {
                       continue
                   }
                   newMask := mask | bit
                   // count inversions added: previous words > i
                   invAdd := bits.OnesCount(uint(mask & highMask[i]))
                   nx := x + invAdd
                   idx := newMask*invSize + nx
                   np := int32(list[j])
                   if dp[idx] > np {
                       dp[idx] = np
                   }
               }
           }
       }
       // find minimal inversions for fullMask
       bestX := -1
       baseF := fullMask * invSize
       for x := 0; x <= invMax; x++ {
           if dp[baseF+x] < INF {
               bestX = x
               break
           }
       }
       if bestX >= 0 {
           p := invMax - bestX + 1
           if p > bestP {
               bestP = p
               bestIdx = prob
           }
       }
   }
   if bestIdx < 0 {
       fmt.Println("Brand new problem!")
   } else {
       fmt.Println(bestIdx)
       bar := strings.Repeat("|", bestP)
       fmt.Println(":" + bar + ":")
   }
}
