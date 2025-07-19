package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   mVal int
   fArr []int
   adjArr []uint64
   tempMax int
)

// dfs explores choices from index x, current size sz, and bitmask cur
func dfs(x, sz int, cur uint64) {
   if x >= mVal {
       if sz > tempMax {
           tempMax = sz
       }
       return
   }
   // pruning: even taking all remaining cannot exceed tempMax
   if sz + fArr[x] <= tempMax {
       return
   }
   // try including x if no conflict
   if adjArr[x] & cur == 0 {
       dfs(x+1, sz+1, cur | (uint64(1) << x))
   }
   // skip x
   dfs(x+1, sz, cur)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   // read n and m
   _, _ = fmt.Fscan(reader, &n, &mVal)
   // mapping of string to id
   idMap := make(map[string]int)
   cnt := 0
   adjArr = make([]uint64, mVal)
   // read operations
   for i := 0; i < n; {
       var tp int
       _, _ = fmt.Fscan(reader, &tp)
       i++
       if tp == 1 {
           continue
       }
       // tp == 2: start a clause
       var s string
       _, _ = fmt.Fscan(reader, &s)
       if _, ok := idMap[s]; !ok {
           cnt++
           idMap[s] = cnt
       }
       aj := uint64(1) << (idMap[s] - 1)
       // read until tp==1 or end
       for i < n {
           var tp2 int
           _, _ = fmt.Fscan(reader, &tp2)
           i++
           if tp2 == 1 {
               break
           }
           _, _ = fmt.Fscan(reader, &s)
           if _, ok := idMap[s]; !ok {
               cnt++
               idMap[s] = cnt
           }
           aj |= uint64(1) << (idMap[s] - 1)
       }
       // update adjacency: for each element in aj, mark edges
       for j := 0; j < mVal; j++ {
           if (aj & (uint64(1) << j)) != 0 {
               adjArr[j] |= aj
           }
       }
   }
   // DP array
   fArr = make([]int, mVal+1)
   // compute from back to front
   for i := mVal - 1; i >= 0; i-- {
       tempMax = 0
       // start with choosing i
       dfs(i+1, 1, uint64(1)<<i)
       // best is max(skip or include)
       if fArr[i+1] > tempMax {
           fArr[i] = fArr[i+1]
       } else {
           fArr[i] = tempMax
       }
   }
   // result is fArr[0]
   fmt.Println(fArr[0])
}
