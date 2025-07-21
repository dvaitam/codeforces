package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read string s
   var s string
   fmt.Fscan(reader, &s)
   n := len(s)
   // Build occurrences
   occ := make([][]int, 26)
   for i := 0; i < n; i++ {
       c := s[i] - 'a'
       occ[c] = append(occ[c], i)
   }
   // Read queries
   var m int
   fmt.Fscan(reader, &m)
   qMasks := make([]uint32, m)
   maskCount := make(map[uint32]int)
   for i := 0; i < m; i++ {
       var ci string
       fmt.Fscan(reader, &ci)
       var mask uint32
       for j := 0; j < len(ci); j++ {
           mask |= 1 << (ci[j] - 'a')
       }
       qMasks[i] = mask
       // initialize count entry
       maskCount[mask] = 0
   }
   // ptr for next occurrence index per letter
   ptr := make([]int, 26)
   // Enumerate windows
   for l := 0; l < n; l++ {
       if l > 0 {
           // advance ptr for letter leaving l-1
           c0 := s[l-1] - 'a'
           if ptr[c0] < len(occ[c0]) && occ[c0][ptr[c0]] == l-1 {
               ptr[c0]++
           }
       }
       var mask uint32
       // process new letters
       for {
           minPos := n
           var minC int = -1
           // find next new letter occurrence
           for c := 0; c < 26; c++ {
               if mask&(1<<c) != 0 {
                   continue
               }
               pi := ptr[c]
               if pi < len(occ[c]) {
                   p := occ[c][pi]
                   if p < minPos {
                       minPos = p
                       minC = c
                   }
               }
           }
           if minC < 0 {
               break
           }
           // window [l, minPos-1] has current mask
           if mask != 0 {
               // check left maximal: either l==0 or s[l-1] not in mask
               if l == 0 || (mask&(1<<uint(s[l-1]-'a')))==0 {
                   if _, ok := maskCount[mask]; ok {
                       maskCount[mask]++
                   }
               }
           }
           // add new letter
           mask |= 1 << uint(minC)
       }
       // final window [l, n-1]
       if mask != 0 {
           if l == 0 || (mask&(1<<uint(s[l-1]-'a')))==0 {
               if _, ok := maskCount[mask]; ok {
                   maskCount[mask]++
               }
           }
       }
   }
   // Output answers
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 0; i < m; i++ {
       fmt.Fprintln(writer, maskCount[qMasks[i]])
   }
}
