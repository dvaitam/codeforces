package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strconv"
   "strings"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read costs
   line, err := reader.ReadString('\n')
   if err != nil && err != io.EOF {
       panic(err)
   }
   fields := strings.Fields(line)
   ti, _ := strconv.Atoi(fields[0])
   td, _ := strconv.Atoi(fields[1])
   tr, _ := strconv.Atoi(fields[2])
   te, _ := strconv.Atoi(fields[3])
   // read strings
   a, _ := reader.ReadString('\n')
   b, _ := reader.ReadString('\n')
   a = strings.TrimSpace(a)
   b = strings.TrimSpace(b)
   n, m := len(a), len(b)
   // maximum distance sentinel
   maxCost := ti
   if td > maxCost {
       maxCost = td
   }
   if tr > maxCost {
       maxCost = tr
   }
   if te > maxCost {
       maxCost = te
   }
   maxdist := (n + m) * maxCost
   // dp matrix H size (n+2)x(m+2)
   H := make([][]int, n+2)
   for i := range H {
       H[i] = make([]int, m+2)
   }
   // initialize
   H[0][0] = maxdist
   for i := 0; i <= n; i++ {
       H[i+1][0] = maxdist
       H[i+1][1] = i * td
   }
   for j := 0; j <= m; j++ {
       H[0][j+1] = maxdist
       H[1][j+1] = j * ti
   }
   // last occurrences for chars in a
   da := make([]int, 26)
   for i := 1; i <= n; i++ {
       db := 0
       for j := 1; j <= m; j++ {
           i1 := da[b[j-1]-'a']
           j1 := db
           cost := tr
           if a[i-1] == b[j-1] {
               cost = 0
               db = j
           }
           // substitution, insertion, deletion
           h0 := H[i][j] + cost
           h1 := H[i+1][j] + ti
           h2 := H[i][j+1] + td
           h := h0
           if h1 < h {
               h = h1
           }
           if h2 < h {
               h = h2
           }
           // transposition
           // H[i1][j1] + (i-i1-1)*td + te + (j-j1-1)*ti
           trans := H[i1][j1]
           trans += (i - i1 - 1) * td
           trans += te
           trans += (j - j1 - 1) * ti
           if trans < h {
               h = trans
           }
           H[i+1][j+1] = h
       }
       da[a[i-1]-'a'] = i
   }
   // result
   res := H[n+1][m+1]
   fmt.Println(res)
}
