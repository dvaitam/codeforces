package main

import (
   "bufio"
   "os"
   "strconv"
)

type Edge struct {
   to, id int
}

func main() {
   r := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   // fast integer reading
   readInt := func() int {
       var x int
       var c byte
       var err error
       // skip non-number
       for {
           c, err = r.ReadByte()
           if err != nil {
               return x
           }
           if (c >= '0' && c <= '9') || c == '-' {
               break
           }
       }
       sign := 1
       if c == '-' {
           sign = -1
           c, _ = r.ReadByte()
       }
       for c >= '0' && c <= '9' {
           x = x*10 + int(c-'0')
           c, _ = r.ReadByte()
       }
       return x * sign
   }

   t := readInt()
   for t > 0 {
       t--
       n := readInt()
       // build tree
       adj := make([][]Edge, n+1)
       deg := make([]int, n+1)
       for i := 1; i < n; i++ {
           u := readInt()
           v := readInt()
           adj[u] = append(adj[u], Edge{v, i})
           adj[v] = append(adj[v], Edge{u, i})
           deg[u]++
           deg[v]++
       }
       // check validity
       ok := true
       for i := 1; i <= n; i++ {
           if deg[i] > 2 {
               ok = false
               break
           }
       }
       if !ok {
           w.WriteString("-1\n")
           continue
       }
       // find leaf
       st := 1
       for i := 1; i <= n; i++ {
           if deg[i] == 1 {
               st = i
               break
           }
       }
       // assign weights along the path
       ans := make([]int, n)
       prev, curr, parity := 0, st, 0
       for {
           found := false
           for _, e := range adj[curr] {
               if e.to == prev {
                   continue
               }
               if parity == 0 {
                   ans[e.id] = 2
               } else {
                   ans[e.id] = 3
               }
               parity ^= 1
               prev = curr
               curr = e.to
               found = true
               break
           }
           if !found {
               break
           }
       }
       // output
       for i := 1; i < n; i++ {
           if i > 1 {
               w.WriteByte(' ')
           }
           w.WriteString(strconv.Itoa(ans[i]))
       }
       w.WriteByte('\n')
   }
}
