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
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       solve(reader, writer)
   }
}

func solve(r *bufio.Reader, w *bufio.Writer) {
   var n int
   fmt.Fscan(r, &n)
   var str string
   fmt.Fscan(r, &str)
   s := []byte(str)
   // Map directions to their indices
   mp := make(map[byte][]int)
   for i, c := range s {
       mp[c] = append(mp[c], i)
   }
   x := len(mp['E']) - len(mp['W'])
   y := len(mp['N']) - len(mp['S'])
   // If differences are odd, impossible
   if x%2 != 0 || y%2 != 0 {
       fmt.Fprintln(w, "NO")
       return
   }
   // Initialize all as 'R'
   ans := make([]byte, n)
   for i := range ans {
       ans[i] = 'R'
   }
   if x != 0 || y != 0 {
       if x > 0 {
           for i := 0; i < x/2; i++ {
               ans[mp['E'][i]] = 'H'
           }
       } else {
           for i := 0; i < (-x)/2; i++ {
               ans[mp['W'][i]] = 'H'
           }
       }
       if y > 0 {
           for i := 0; i < y/2; i++ {
               ans[mp['N'][i]] = 'H'
           }
       } else {
           for i := 0; i < (-y)/2; i++ {
               ans[mp['S'][i]] = 'H'
           }
       }
   } else {
       // Balanced but no excess; need at least one H-pair
       if n == 2 {
           fmt.Fprintln(w, "NO")
           return
       }
       if len(mp['E']) > 0 {
           ans[mp['E'][0]] = 'H'
           ans[mp['W'][0]] = 'H'
       } else {
           ans[mp['N'][0]] = 'H'
           ans[mp['S'][0]] = 'H'
       }
   }
   fmt.Fprintln(w, string(ans))
}
