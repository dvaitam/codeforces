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

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n+1)
       for i := 1; i <= n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       ans := make([]int, n+1)
       vis := make([]bool, n+1)
       // first, assign as many wishes as possible
       matched := 0
       for i := 1; i <= n; i++ {
           if !vis[a[i]] {
               vis[a[i]] = true
               ans[i] = a[i]
               matched++
           }
       }
       // output number of matches
       fmt.Fprintln(writer, matched)
       // collect unused targets, from n down to 1 so pop gives smallest
       st := make([]int, 0, n)
       for i := n; i >= 1; i-- {
           if !vis[i] {
               st = append(st, i)
           }
       }
       // fill remaining
       for i := 1; i <= n; i++ {
           if ans[i] == 0 {
               // pop last
               last := st[len(st)-1]
               st = st[:len(st)-1]
               ans[i] = last
           }
       }
       // fix self-assignments
       las := 0
       for i := 1; i <= n; i++ {
           if ans[i] == i {
               if las == 0 {
                   las = i
               } else {
                   ans[i], ans[las] = ans[las], ans[i]
               }
           }
       }
       if las != 0 {
           // find someone with same original wish as las
           for i := 1; i <= n; i++ {
               if a[i] == a[las] {
                   ans[i], ans[las] = ans[las], ans[i]
                   break
               }
           }
       }
       // print assignment
       for i := 1; i <= n; i++ {
           if i > 1 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, ans[i])
       }
       writer.WriteByte('\n')
   }
}
