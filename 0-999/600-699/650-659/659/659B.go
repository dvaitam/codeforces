package main

import (
   "bufio"
   "fmt"
   "os"
)

type student struct {
   name   string
   points int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   fmt.Fscan(in, &n, &m)

   // top[r][0], top[r][1], top[r][2] are first, second, third best
   top := make([][3]student, m+1)
   for i := 1; i <= m; i++ {
       for j := 0; j < 3; j++ {
           top[i][j].points = -1
       }
   }

   for i := 0; i < n; i++ {
       var name string
       var r, p int
       fmt.Fscan(in, &name, &r, &p)
       if p >= top[r][0].points {
           top[r][2] = top[r][1]
           top[r][1] = top[r][0]
           top[r][0].points = p
           top[r][0].name = name
       } else if p >= top[r][1].points {
           top[r][2] = top[r][1]
           top[r][1].points = p
           top[r][1].name = name
       } else if p >= top[r][2].points {
           top[r][2].points = p
           top[r][2].name = name
       }
   }

   for i := 1; i <= m; i++ {
       if top[i][1].points == top[i][2].points {
           fmt.Fprintln(out, "?")
       } else {
           fmt.Fprintf(out, "%s %s\n", top[i][0].name, top[i][1].name)
       }
   }
}
