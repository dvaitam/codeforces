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
template<typename T> il T chkmax(T &a, const T &b) {return a = (b > a ? b : a);}
template<typename T> il T chkmin(T &a, const T &b) {return a = (b < a ? b : a);}
template<typename T> il T myabs(const T &a) {return a >= 0 ? a : -a;}
template<typename T> il void myswap(T &a, T &b) {
    T t = a;
    a = b, b = t;
    return;
}

const int maxn = 2e5 + 5;
vector<int> G0[maxn], G[maxn];
int n, m, vis[maxn], col[maxn], ind[maxn], x[maxn];
bool flg = 1;

struct Relation {
    int op, u, v;
} a[maxn];

void dfs(int u, int cur) {
    col[u] = cur, vis[u] = 1;
    for (int &v : G0[u]) {
        if (vis[v] && col[v] == cur) flg = 0;
        if (vis[v]) continue;
        dfs(v, cur == 1 ? 2 : 1);
    }
}

int main() {
    read(n, m);
    FOR(i, 1, m) read(a[i].op, a[i].u, a[i].v), G0[a[i].u].push_back(a[i].v), G0[a[i].v].push_back(a[i].u);
    FOR(i, 1, n) if (!vis[i]) dfs(i, 1);
    if (!flg) {
        print("NO");
    } else {
        FOR(i, 1, m) {
            const int &u = a[i].u, &v = a[i].v;
            if (a[i].op == 1) {
                if (col[u] == 1) G[u].push_back(v), ++ind[v];
                else G[v].push_back(u), ++ind[u];
            } else {
                if (col[u] == 2) G[u].push_back(v), ++ind[v];
                else G[v].push_back(u), ++ind[u];
            }
        }
        queue<int> q;
        FOR(i, 1, n) if (!ind[i]) q.push(i);
        int cntx = 0;
        while (!q.empty()) {
            int u = q.front();
            x[u] = ++cntx;
            q.pop();
            for (const int &v : G[u]) {
                if (!--ind[v]) q.push(v);
            }
        }
        if (cntx != n) {
            print("NO");
        } else {
            print("YES");
            FOR(i, 1, n) {
                putchar(col[i] == 1 ? 'L' : 'R'); putchar(' ');
                print(x[i]);
            }
        }
    }
    return output(), 0;
}

} // namespace YangTY

int main() {
    YangTY::main();
    return 0;
}