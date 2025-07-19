package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   t := readInt(reader)
   for ; t > 0; t-- {
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   n := readInt(reader)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       a[i] = readInt(reader) - 1
   }
   freq := make([]int, n)
   maxFreq := 0
   for i := 0; i < n; i++ {
       freq[a[i]]++
       if freq[a[i]] > maxFreq {
           maxFreq = freq[a[i]]
       }
   }
   pos := make([][]int, n)
   for i := 0; i < n; i++ {
       v := a[i]
       pos[v] = append(pos[v], i)
   }
   sort.Slice(pos, func(i, j int) bool {
       return len(pos[i]) > len(pos[j])
   })
   res := make([]int, n)
   type pair struct { val, idx int }
   for rep := 0; rep < maxFreq; rep++ {
       var cyc []pair
       for i := 0; i < len(pos); i++ {
           if len(pos[i]) == 0 {
               break
           }
           idx := pos[i][len(pos[i])-1]
           pos[i] = pos[i][:len(pos[i])-1]
           cyc = append(cyc, pair{a[idx], idx})
       }
       sort.Slice(cyc, func(i, j int) bool {
           return cyc[i].val < cyc[j].val
       })
       m := len(cyc)
       for i := 0; i < m; i++ {
           next := (i + 1) % m
           res[cyc[i].idx] = cyc[next].idx
       }
   }
   for i := 0; i < n; i++ {
       writer.WriteString(fmt.Sprintf("%d ", a[res[i]]+1))
   }
   writer.WriteString("\n")
}

func readInt(r *bufio.Reader) int {
   var c byte
   var err error
   c, err = r.ReadByte()
   if err != nil {
       return 0
   }
   for (c < '0' || c > '9') && c != '-' {
       c, err = r.ReadByte()
       if err != nil {
           return 0
       }
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = r.ReadByte()
   }
   x := 0
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, err = r.ReadByte()
       if err != nil {
           break
       }
   }
   return x * sign
}
