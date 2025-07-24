package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strconv"
   "strings"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var line string
   if _, err := fmt.Fscanln(in, &line); err != nil {
       return
   }
   n, _ := strconv.Atoi(line)
   pre := make([]int, n)
   cyc := make([]int, n)
   knownCycles := make(map[int]bool)
   for i := 0; i < n; i++ {
       var a, b string
       fmt.Fscan(in, &a, &b)
       if a == "?" {
           pre[i] = -1
       } else {
           v, _ := strconv.Atoi(a)
           pre[i] = v
       }
       if b == "?" {
           cyc[i] = -1
       } else {
           v, _ := strconv.Atoi(b)
           cyc[i] = v
           knownCycles[v] = true
       }
   }
   // prepare output and flags
   f := make([]int, n)
   for i := range f {
       f[i] = -1
   }
   fSet := make([]bool, n)

   // process known cycle lengths in increasing order
   var lengths []int
   for L := range knownCycles {
       lengths = append(lengths, L)
   }
   sort.Ints(lengths)
   for _, L := range lengths {
       // collect nodes with cyc == L
       var P0 []int
       var SU []int
       depthMap := make(map[int][]int)
       for i := 0; i < n; i++ {
           if cyc[i] != L {
               continue
           }
           if pre[i] == 0 {
               P0 = append(P0, i)
           } else if pre[i] > 0 {
               depthMap[pre[i]] = append(depthMap[pre[i]], i)
           } else {
               // pre unknown
               SU = append(SU, i)
           }
       }
       if len(P0) > L {
           fmt.Println(-1)
           return
       }
       // select extra cycle nodes
       need := L - len(P0)
       var cycleNodes []int
       cycleNodes = append(cycleNodes, P0...)
       // take from SU first
       take := min(need, len(SU))
       cycleNodes = append(cycleNodes, SU[:take]...)
       // mark these SU consumed
       SU = SU[take:]
       need -= take
       // take remaining from global unassigned candidates
       if need > 0 {
           for i := 0; i < n && need > 0; i++ {
               if cyc[i] != -1 || pre[i] > 0 {
                   continue
               }
               // candidate: pre[i] == 0 or unknown, cyc unknown
               if pre[i] == 0 || pre[i] == -1 {
                   // assign to this cycle
                   cyc[i] = L
                   pre[i] = 0
                   cycleNodes = append(cycleNodes, i)
                   need--
               }
           }
           if need > 0 {
               fmt.Println(-1)
               return
           }
       }
       // assign cycle edges
       if len(cycleNodes) != L {
           fmt.Println(-1)
           return
       }
       // order arbitrary
       for i, v := range cycleNodes {
           next := cycleNodes[(i+1)%L]
           f[v] = next
           fSet[v] = true
           pre[v] = 0
           cyc[v] = L
       }
       // remaining SU become depth 1
       for _, v := range SU {
           pre[v] = 1
       }
       // include SU and known deep nodes into depthMap
       // SU now all depth 1
       if len(SU) > 0 {
           depthMap[1] = append(depthMap[1], SU...)
       }
       // also include any newly assigned from global in earlier selection: they have pre=0, so in cycleNodes
       // build depths up to max
       maxd := 0
       for d := range depthMap {
           if d > maxd {
               maxd = d
           }
       }
       // for depths k=1..maxd
       for d := 1; d <= maxd; d++ {
           prev := []int{}
           if d-1 == 0 {
               prev = cycleNodes
           } else {
               prev = depthMap[d-1]
           }
           if len(depthMap[d]) > 0 && len(prev) == 0 {
               fmt.Println(-1)
               return
           }
           for _, v := range depthMap[d] {
               // assign to first prev
               f[v] = prev[0]
               fSet[v] = true
               if cyc[v] != L {
                   cyc[v] = L
               }
           }
       }
   }
   // handle remaining unassigned
   var rem []int
   for i := 0; i < n; i++ {
       if !fSet[i] {
           rem = append(rem, i)
       }
   }
   if len(rem) > 0 {
       // find S0
       var S0 []int
       for _, v := range rem {
           if pre[v] == 0 || pre[v] == -1 {
               S0 = append(S0, v)
           }
       }
       if len(S0) == 0 {
           fmt.Println(-1)
           return
       }
       // create cycle length 1
       root := S0[0]
       pre[root] = 0
       cyc[root] = 1
       f[root] = root
       fSet[root] = true
       // other S0 become depth1
       for i := 1; i < len(S0); i++ {
           v := S0[i]
           if pre[v] == 0 {
               fmt.Println(-1)
               return
           }
           pre[v] = 1
           cyc[v] = 1
           f[v] = root
           fSet[v] = true
       }
       // other rem with pre>0
       for _, v := range rem {
           if fSet[v] {
               continue
           }
           // pre[v] >0 assumed
           if pre[v] <= 0 {
               pre[v] = 1
           }
           cyc[v] = 1
           // ensure depthMap includes depth1 if needed
           // direct attach
           f[v] = root
           fSet[v] = true
       }
   }
   // output
   w := bufio.NewWriter(os.Stdout)
   for i := 0; i < n; i++ {
       if f[i] < 0 {
           fmt.Println(-1)
           return
       }
       w.WriteString(strconv.Itoa(f[i]+1))
       if i+1 < n {
           w.WriteByte(' ')
       }
   }
   w.WriteByte('\n')
   w.Flush()
}

func min(a, b int) int { if a < b { return a }; return b }
