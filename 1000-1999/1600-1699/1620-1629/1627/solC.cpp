#include <bits/stdc++.h>
#define il inline
#define FOR(i, a, b) for (int i = (a); i <= (b); ++i)
#define DEC(i, a, b) for (int i = (a); i >= (b); --i)

using namespace std;

namespace YangTY {
namespace fastIO {
const int maxc = 1 << 23;
char ibuf[maxc], *__p1 = ibuf, *__p2 = ibuf;
il char getchar() {return __p1 == __p2 && (__p2 = (__p1 = ibuf) + fread(ibuf, 1, maxc, stdin), __p1 == __p2) ? EOF : *__p1++;}
template<typename T> void read(T &n) {
    int x = 0; n = 0;
    char c = getchar();
    while (!isdigit(c))
        x |= (c == '-'), c = getchar();
    while (isdigit(c))
        n = n * 10 + c - '0', c = getchar();
    n = x ? -n : n;
}
void read(char *s) {
    int p = 0;
    char c = getchar();
    while (isspace(c)) c = getchar();
    while (~c && !isspace(c)) s[p++] = c, c = getchar();
    return;
}
template<typename T1, typename... T2> void read(T1 &a, T2&... x) {
    read(a), read(x...);
    return;
}
char obuf[maxc], *__pO = obuf;
il void putchar(char c) {*__pO++ = c;}
template<typename T> void print(T x, char c = '\n') {
    static char stk[50];
    int top = 0;
    if (x < 0) putchar('-'), x = -x;
    if (x) {
        for (; x; x /= 10) stk[++top] = x % 10 + '0';
        while (top) putchar(stk[top--]);
    } else putchar('0');
    putchar(c);
    return;
}
void print(char *s, char c = '\n') {
    for (int i = 0; s[i]; ++i) putchar(s[i]);
    putchar(c);
    return;
}
void print(const char *s, char c = '\n') {
    for (int i = 0; s[i]; ++i) putchar(s[i]);
    putchar(c);
    return;
}
template<typename T1, typename... T2> il void print(T1 a, T2... x) {
    if (!sizeof...(x)) print(a);
    else print(a, ' '), print(x...);
    return;
}
void output() {fwrite(obuf, __pO - obuf, 1, stdout);}
} // namespace fastIO

using namespace fastIO;

template<typename T> il T max(const T &a, const T &b) {return a > b ? a : b;}
template<typename T> il T min(const T &a, const T &b) {return a < b ? a : b;}
template<typename T, typename...Args> il T max(const T &a, const Args&... args) {
    T res = max(args...);
    return max(a, res);
}
template<typename T, typename...Args> il T min(const T &a, const Args&... args) {
    T res = min(args...);
    return min(a, res);
}
template<typename T> il T chkmax(T &a, const T &b) {return a = max(a, b);}
template<typename T> il T chkmin(T &a, const T &b) {return a = min(a, b);}
template<typename T> il T myabs(const T &a) {return a >= 0 ? a : -a;}
template<typename T> il void myswap(T &a, T &b) {
    T t = a;
    a = b, b = t;
    return;
}

const int maxn = 1e5 + 5, N = 1e5;
int isp[maxn], pri[maxn], tot;

void sieve() {
    FOR(i, 2, N) {
        if (!isp[i]) pri[++tot] = i;
        for (int j = 1; j <= tot && i * pri[j] <= N; ++j) {
            isp[i * pri[j]] = 1;
            if (i % pri[j] == 0) break;
        }
    }
    return;
}

struct edge {
    int to, nxt, id;
} e[maxn << 1];
int head[maxn], cnte, deg[maxn], n, ans[maxn];

void add(int u, int v, int id) {
    e[++cnte].to = v;
    e[cnte].nxt = head[u];
    e[cnte].id = id;
    head[u] = cnte;
    return;
}

void dfs(int u, int fa, int cur) {
    for (int i = head[u]; i; i = e[i].nxt) {
        if (e[i].to == fa) continue;
        ans[e[i].id] = cur ? 2 : 3;
        dfs(e[i].to, u, cur ^ 1);
    }
    return;
}

int main() {
    sieve();
    int T; read(T);
    while (T--) {
        read(n);
        FOR(i, 1, n) head[i] = deg[i] = 0;
        cnte = 0;
        FOR(i, 1, n - 1) {
            int u, v; read(u, v);
            ++deg[u], ++deg[v];
            add(u, v, i), add(v, u, i);
        }
        bool flg = 1;
        FOR(i, 1, n) if (deg[i] >= 3) flg = 0;
        if (!flg) print(-1);
        else {
            int st = 0;
            FOR(i, 1, n) if (deg[i] == 1) {
                st = i;
                break;
            }
            dfs(st, 0, 0);
            FOR(i, 1, n - 1) print(ans[i], ' ');
            putchar('\n');
        }
    }
    return output(), 0;
}

} // namespace YangTY

int main() {
    YangTY::main();
    return 0;
}