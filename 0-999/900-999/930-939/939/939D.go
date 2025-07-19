package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   var s1, s2 string
   fmt.Fscan(reader, &s1, &s2)

   parent := make([]int, 27)
   for i := 1; i <= 26; i++ {
       parent[i] = i
   }
   var ops [][2]int

   var find func(int) int
   find = func(x int) int {
       if parent[x] != x {
           parent[x] = find(parent[x])
       }
       return parent[x]
   }
   union := func(a, b int) {
       pa := find(a)
       pb := find(b)
       if pa != pb {
           parent[pa] = pb
       }
   }

   for i := 0; i < n; i++ {
       a := int(s1[i]-'a'+1)
       b := int(s2[i]-'a'+1)
       if a != b && find(a) != find(b) {
           ops = append(ops, [2]int{a, b})
           union(a, b)
       }
   }

   fmt.Fprintln(writer, len(ops))
   for _, op := range ops {
       a, b := op[0], op[1]
       fmt.Fprintf(writer, "%c %c\n", rune(a-1+'a'), rune(b-1+'a'))
   }
}
