import java.io.*;
import java.util.*;
 
public class Main {
    static final class FastScanner {
        private final InputStream in = System.in;
        private final byte[] buffer = new byte[1 << 16];
        private int ptr = 0, len = 0;
 
        private int read() throws IOException {
            if (ptr >= len) {
                len = in.read(buffer);
                ptr = 0;
                if (len <= 0) return -1;
            }
            return buffer[ptr++];
        }
 
        int nextInt() throws IOException {
            int c;
            do {
                c = read();
            } while (c <= ' ' && c != -1);
 
            int sign = 1;
            if (c == '-') {
                sign = -1;
                c = read();
            }
 
            int val = 0;
            while (c > ' ') {
                val = val * 10 + (c - '0');
                c = read();
            }
            return val * sign;
        }
    }
 
    static final class SegTree {
        int n;
        int size;
        int[] initMax; // initial relative maximums
        int[] max;     // current relative maximums
        int[] lazy;    // lazy add
 
        void prepare(int n) {
            if (this.n == n && initMax != null) return;
            this.n = n;
            this.size = 4 * n + 5;
 
            if (initMax == null || initMax.length < size) {
                initMax = new int[size];
                max = new int[size];
                lazy = new int[size];
            }
 
            buildInit(1, 1, n);
        }
 
        private void buildInit(int node, int l, int r) {
            // Initial array is [0, 1, 2, ..., n-1]
            // So maximum on segment [l, r] is r-1.
            initMax[node] = r - 1;
            if (l == r) return;
            int mid = (l + r) >>> 1;
            buildInit(node << 1, l, mid);
            buildInit(node << 1 | 1, mid + 1, r);
        }
 
        void reset() {
            System.arraycopy(initMax, 0, max, 0, size);
            Arrays.fill(lazy, 0, size, 0);
        }
 
        int rootMax() {
            return max[1];
        }
 
        private void apply(int node, int delta) {
            max[node] += delta;
            lazy[node] += delta;
        }
 
        private void push(int node) {
            int add = lazy[node];
            if (add != 0) {
                apply(node << 1, add);
                apply(node << 1 | 1, add);
                lazy[node] = 0;
            }
        }
 
        void addPrefix(int r) {
            if (r <= 0) return;
            addPrefix(1, 1, n, r);
        }
 
        private void addPrefix(int node, int l, int r, int limit) {
            if (r <= limit) {
                apply(node, 1);
                return;
            }
            if (l == r) return;
 
            push(node);
            int mid = (l + r) >>> 1;
            addPrefix(node << 1, l, mid, limit);
            if (limit > mid) {
                addPrefix(node << 1 | 1, mid + 1, r, limit);
            }
            max[node] = Math.max(max[node << 1], max[node << 1 | 1]);
        }
 
        int firstGE(int target) {
            if (max[1] < target) return n + 1;
            return firstGE(1, 1, n, target);
        }
 
        private int firstGE(int node, int l, int r, int target) {
            if (l == r) return l;
            push(node);
            int mid = (l + r) >>> 1;
            if (max[node << 1] >= target) {
                return firstGE(node << 1, l, mid, target);
            } else {
                return firstGE(node << 1 | 1, mid + 1, r, target);
            }
        }
    }
 
    int n, k, m;
    int[] a;
    int[] leftMax;
    final SegTree seg = new SegTree();
 
    private boolean can(int p) {
        final int n = this.n;
        final int m = this.m;
        final int[] a = this.a;
        final int[] leftMax = this.leftMax;
 
        // We work with h[j] = (j-1) + best_j.
        // Initially: h[j] = j-1, so target for feasibility is m-1.
        final int target = m - 1;
        final int shift = p - m; // need = a[i] - shift = a[i] - p + m
 
        // Left scan: for every position, compute the largest possible center rank.
        seg.reset();
        int border = m; // first j such that h[j] >= m-1
 
        for (int i = 0; i < n; i++) {
            leftMax[i] = m - border + 1;
 
            int need = a[i] - shift;
            if (need > 0) {
                int r;
                if (need > seg.rootMax()) {
                    r = m;
                } else {
                    r = seg.firstGE(need) - 1;
                }
 
                if (r > 0) {
                    seg.addPrefix(r);
                    if (border > 1 && r >= border - 1) {
                        border = seg.firstGE(target);
                    }
                }
            }
        }
 
        // Right scan: for every position, compute the smallest possible center rank.
        seg.reset();
        border = m;
 
        for (int i = n - 1; i >= 0; i--) {
            int rightMin = border;
 
            if (a[i] >= p && rightMin <= leftMax[i]) {
                return true;
            }
 
            int need = a[i] - shift;
            if (need > 0) {
                int r;
                if (need > seg.rootMax()) {
                    r = m;
                } else {
                    r = seg.firstGE(need) - 1;
                }
 
                if (r > 0) {
                    seg.addPrefix(r);
                    if (border > 1 && r >= border - 1) {
                        border = seg.firstGE(target);
                    }
                }
            }
        }
 
        return false;
    }
 
    private void solve() throws Exception {
        FastScanner fs = new FastScanner();
        StringBuilder out = new StringBuilder();
 
        int t = fs.nextInt();
        while (t-- > 0) {
            n = fs.nextInt();
            k = fs.nextInt();
            m = n - k;
 
            if (a == null || a.length < n) a = new int[n];
            if (leftMax == null || leftMax.length < n) leftMax = new int[n];
 
            int maxA = 0;
            for (int i = 0; i < n; i++) {
                a[i] = fs.nextInt();
                if (a[i] > maxA) maxA = a[i];
            }
 
            seg.prepare(m);
 
            int lo = 1, hi = maxA, ans = 1;
            while (lo <= hi) {
                int mid = (lo + hi) >>> 1;
                if (can(mid)) {
                    ans = mid;
                    lo = mid + 1;
                } else {
                    hi = mid - 1;
                }
            }
 
            out.append(ans).append('\n');
        }
 
        System.out.print(out.toString());
    }
 
    public static void main(String[] args) throws Exception {
        new Main().solve();
    }
}