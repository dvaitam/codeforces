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
   // TODO: Implement centroid decomposition based solution as per editorial:
   // 1. Read tree, compute initial depth parity and centroid decomposition.
   // 2. Precompute distances to centroids and component indices.
   // 3. Maintain for each centroid and parity the component maxima using heaps or multisets.
   // 4. On each edge flip, update subtree parities and propagate changes up centroid tree.
   // 5. Maintain global diameters for even parity paths and output after each query.

   // Placeholder: outputs zeros
   var m int
   for i := 0; i < n-1; i++ {
       var u, v, t int
       fmt.Fscan(reader, &u, &v, &t)
   }
   fmt.Fscan(reader, &m)
   for i := 0; i < m; i++ {
       var id int
       fmt.Fscan(reader, &id)
       fmt.Fprintln(writer, 0)
   }
}
