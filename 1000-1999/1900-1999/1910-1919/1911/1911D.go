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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
      return
   }

   freq := make(map[int]int, n)
   values := make([]int, n)
   for i := 0; i < n; i++ {
      fmt.Fscan(reader, &values[i])
      freq[values[i]]++
      if freq[values[i]] > 2 {
         fmt.Fprintln(writer, "NO")
         return
      }
   }

   asc := make([]int, 0, len(freq))
   desc := make([]int, 0)
   for v, c := range freq {
      asc = append(asc, v)
      if c == 2 {
         desc = append(desc, v)
      }
   }
   sort.Ints(asc)
   sort.Sort(sort.Reverse(sort.IntSlice(desc)))

   fmt.Fprintln(writer, "YES")
   printSequence(writer, asc)
   printSequence(writer, desc)
}

func printSequence(writer *bufio.Writer, seq []int) {
   fmt.Fprint(writer, len(seq))
   for _, v := range seq {
      fmt.Fprint(writer, " ")
      fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
