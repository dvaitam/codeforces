#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <cmath>
#include <cctype>
#include <climits>
#include <cassert>
#include <ctime>
#include <iostream>
#include <algorithm>
#include <functional>

#define x first
#define y second
#define MP std::make_pair
#define VAL(x) #x " = " << x
#define DEBUG(...) fprintf(stderr, __VA_ARGS__)

typedef long long LL;
typedef std::pair<int, int> Pii;

const int oo = 0x3f3f3f3f;

template<typename T> inline bool chkmax(T &a, T b) { return a < b ? a = b, true : false; }
template<typename T> inline bool chkmin(T &a, T b) { return a > b ? a = b, true : false; }
template<typename T> T read(T &x)
{
    int f = 1;
    char ch = getchar();
    for (; !isdigit(ch); ch = getchar())
        if (ch == '-')
            f = -1;
    for (x = 0; isdigit(ch); ch = getchar())
        x = 10 * x + ch - '0';
    return x *= f;
}
template<typename T> void write(T x)
{
    if (x == 0) {
        putchar('0');
        return;
    }
    if (x < 0) {
        putchar('-');
        x = -x;
    }
    static char s[20];
    int top = 0;
    for (; x; x /= 10)
        s[++top] = x % 10 + '0';
    while (top)
        putchar(s[top--]);
}
// EOT

const int MAXN = 3e5 + 5;

struct Edge {
    int v, next;

    Edge(int v0 = 0, int next0 = 0): v(v0), next(next0) { }
};

int N, Q;
int tote, head[MAXN];
Edge edge[MAXN];
int fa[MAXN], size[MAXN], hson[MAXN], centroid[MAXN];

#define MAX_SUBTREE_SIZE(rt, u) (std::max(size[rt] - size[u], size[hson[u]]))

void dfs(int u)
{
    size[u] = 1;
    for (int i = head[u]; i; i = edge[i].next) {
        int v = edge[i].v;
        fa[v] = u;
        dfs(v);
        size[u] += size[v];
        if (size[hson[u]] < size[v])
            hson[u] = v;
    }
    centroid[u] = u;
    if (hson[u] != 0) {
        int p = centroid[hson[u]];
        while (fa[p] != u && MAX_SUBTREE_SIZE(u, fa[p]) <= MAX_SUBTREE_SIZE(u, p)) {
            p = fa[p];
        }
        if (size[hson[u]] > MAX_SUBTREE_SIZE(u, p))
            centroid[u] = p;
    }
}

inline void addEdge(int u, int v)
{
    edge[++tote] = Edge(v, head[u]);
    head[u] = tote;
}

void input()
{
    read(N); read(Q);
    for (int i = 2; i <= N; ++i) {
        int f;
        read(f);
        addEdge(f, i);
    }
}

void solve()
{
    dfs(1);
    while (Q--) {
        int u;
        read(u);
        write(centroid[u]); putchar('\n');
    }
}

int main()
{
#ifndef ONLINE_JUDGE
    freopen("tmp.in", "r", stdin);
    freopen("tmp.out", "w", stdout);
#endif

    input();
    solve();

    return 0;
}