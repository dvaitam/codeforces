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

   var m, k int
   fmt.Fscan(reader, &m, &k)
   fav := make([]bool, m+1)
   for i := 0; i < k; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x >= 1 && x <= m {
           fav[x] = true
       }
   }
   var n int
   fmt.Fscan(reader, &n)
   minFav := make([]int, n)
   maxFav := make([]int, n)
   for i := 0; i < n; i++ {
       var title string
       var d int
       fmt.Fscan(reader, &title)
       fmt.Fscan(reader, &d)
       knownFav := 0
       knownNonFav := 0
       zeros := 0
       for j := 0; j < d; j++ {
           var b int
           fmt.Fscan(reader, &b)
           if b == 0 {
               zeros++
           } else if fav[b] {
               knownFav++
           } else {
               knownNonFav++
           }
       }
       availableNonFav := (m - k) - knownNonFav
       extraMin := 0
       if zeros > availableNonFav {
           extraMin = zeros - availableNonFav
       }
       minFav[i] = knownFav + extraMin
       availableFav := k - knownFav
       extraMax := zeros
       if extraMax > availableFav {
           extraMax = availableFav
       }
       maxFav[i] = knownFav + extraMax
   }

   for i := 0; i < n; i++ {
       isSureFav := true
       isSureNotFav := false
       for j := 0; j < n; j++ {
           if i == j {
               continue
           }
           if maxFav[j] > minFav[i] {
               isSureFav = false
           }
           if minFav[j] > maxFav[i] {
               isSureNotFav = true
           }
       }
       if isSureFav {
           fmt.Fprintln(writer, 0)
       } else if isSureNotFav {
           fmt.Fprintln(writer, 1)
       } else {
           fmt.Fprintln(writer, 2)
       }
   }
}
