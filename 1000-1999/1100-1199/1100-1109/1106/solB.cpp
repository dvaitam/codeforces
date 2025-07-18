#include <cctype>
#include <cstdio>
#include <climits>
#include <algorithm>

template <typename T> inline void read(T& x) {
    int f = 0, c = getchar(); x = 0;
    while (!isdigit(c)) f |= c == '-', c = getchar();
    while (isdigit(c)) x = x * 10 + c - 48, c = getchar();
    if (f) x = -x;
}
template <typename T, typename... Args>
inline void read(T& x, Args&... args) {
    read(x); read(args...); 
}
template <typename T> void write(T x) {
    if (x < 0) x = -x, putchar('-');
    if (x > 9) write(x / 10);
    putchar(x % 10 + 48);
}
template <typename T> void writeln(T x) { write(x); puts(""); }
template <typename T> inline bool chkmin(T& x, const T& y) { return y < x ? (x = y, true) : false; }
template <typename T> inline bool chkmax(T& x, const T& y) { return x < y ? (x = y, true) : false; }

const int maxn = 1e5 + 207;
struct Dish {
    int a, c, id;
    Dish(int x, int y) : a(x), c(y) {}
    Dish() : a(0), c(0) {}
};
Dish a[maxn];
int n, m, pos[maxn], p;
inline bool operator<(const Dish &lhs, const Dish &rhs) {
    return lhs.c < rhs.c || (lhs.c == rhs.c && lhs.id < rhs.id);
}
int main() {
    read(n, m);
    for (int i = 1; i <= n; ++i) read(a[i].a);
    for (int i = 1; i <= n; ++i) read(a[i].c), a[i].id = i;
    std::sort(a + 1, a + n + 1);
    for (int i = 1; i <= n; ++i)
        pos[a[i].id] = i;
    p = 1;
    while (m--) {
        int t, d; read(t, d);
        t = pos[t];
        if (a[t].a >= d) {
            a[t].a -= d;
            writeln(1ll * d * a[t].c);
        } else {
            long long ans = 1ll * a[t].a * a[t].c;
            d -= a[t].a;
            a[t].a = 0;
            while (0207) {
                while (!a[p].a && p < n) ++p;
                if (p == n && !a[p].a) break;
                if (a[p].a >= d) {
                    ans += 1ll * d * a[p].c;
                    a[p].a -= d;
                    d = 0;
                    break;
                } else {
                    ans += 1ll * a[p].a * a[p].c;
                    d -= a[p].a;
                    a[p].a = 0;
                }
            }
            if (d) puts("0");
            else writeln(ans);
        }
    }
    return 0;
}