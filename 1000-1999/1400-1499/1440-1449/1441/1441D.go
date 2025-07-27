package main

import (
   "bufio"
   "fmt"
   "os"
)

// Fast IO
var rdr = bufio.NewReader(os.Stdin)
var wrtr = bufio.NewWriter(os.Stdout)

func readInt() (int, error) {
   var x int
   var neg bool
   b := make([]byte, 0, 16)
   for {
       ch, err := rdr.ReadByte()
       if err != nil {
           return 0, err
       }
       if (ch >= '0' && ch <= '9') || ch == '-' {
           b = append(b, ch)
           break
       }
   }
   for {
       ch, err := rdr.ReadByte()
       if err != nil {
           break
       }
       if ch < '0' || ch > '9' {
           break
       }
       b = append(b, ch)
   }
   // parse b
   for i, ch := range b {
       if ch == '-' && i == 0 {
           neg = true
       } else {
           x = x*10 + int(ch-'0')
       }
   }
   if neg {
       x = -x
   }
   return x, nil
}

func main() {
   defer wrtr.Flush()
   t, _ := readInt()
   for tc := 0; tc < t; tc++ {
       n, _ := readInt()
       a := make([]int, n)
       cntW, cntB := 0, 0
       for i := 0; i < n; i++ {
           ai, _ := readInt()
           a[i] = ai
           if ai == 1 {
               cntW++
           } else if ai == 2 {
               cntB++
           }
       }
       // build adjacency
       adj := make([][]int, n)
       for i := 0; i < n-1; i++ {
           u, _ := readInt(); v, _ := readInt()
           u--; v--
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // if only one color or none, can remove all in one op
       if cntW == 0 || cntB == 0 {
           fmt.Fprintln(wrtr, 1)
           continue
       }
       // count components on white+grey (a != 2)
       vis := make([]bool, n)
       compWG := 0
       stack := make([]int, 0, n)
       for i := 0; i < n; i++ {
           if !vis[i] && a[i] != 2 {
               compWG++
               // dfs
               stack = stack[:0]
               stack = append(stack, i)
               vis[i] = true
               for len(stack) > 0 {
                   v := stack[len(stack)-1]
                   stack = stack[:len(stack)-1]
                   for _, u := range adj[v] {
                       if !vis[u] && a[u] != 2 {
                           vis[u] = true
                           stack = append(stack, u)
                       }
                   }
               }
           }
       }
       // count components on black+grey (a != 1)
       vis2 := make([]bool, n)
       compBG := 0
       for i := 0; i < n; i++ {
           if !vis2[i] && a[i] != 1 {
               compBG++
               stack = stack[:0]
               stack = append(stack, i)
               vis2[i] = true
               for len(stack) > 0 {
                   v := stack[len(stack)-1]
                   stack = stack[:len(stack)-1]
                   for _, u := range adj[v] {
                       if !vis2[u] && a[u] != 1 {
                           vis2[u] = true
                           stack = append(stack, u)
                       }
                   }
               }
           }
       }
       // answer
       // both groups non-empty
       // minimum ops = 1 + min(compWG, compBG)
       ans := 1 + compWG
       if 1+compBG < ans {
           ans = 1 + compBG
       }
       fmt.Fprintln(wrtr, ans)
   }
}
