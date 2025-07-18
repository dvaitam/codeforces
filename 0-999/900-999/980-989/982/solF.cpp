#include <bits/stdc++.h>
#define N 100005
#define M 500005
using namespace std;
int n, m, ans, res, flag[N], tag[N], check[N], sum[N];
int k, la[N], ff[M], c[N], tot, p[N], pre[N], pos[N];
struct node {
    int a, b;
} e[M];

void add(int a, int b) {
    e[++k] = (node){a, b};
    ff[k] = la[a];
    la[a] = k;
}

// find a cycle and store it in p in backwards order
void dfs1(int x) {
    if (tot || check[x])
        return;
    c[x] = -1;
    for (int a = la[x]; a && !tot; a = ff[a]) {
        if (check[e[a].b])
            continue;
        if (c[e[a].b] == -1) {
            while (x != e[a].b)
                p[++tot] = x, x = pre[x];
            p[++tot] = e[a].b;
            return;
        }
        if (!c[e[a].b])
            pre[e[a].b] = x, dfs1(e[a].b);
    }
    c[x] = 1;
}

void dfs2(int S, int x) {
    if (check[x])
        return;
    flag[x] = 1;
    for (int a = la[x]; a; a = ff[a]) {
        if (check[e[a].b])
            continue;
        if (tag[e[a].b]) {
            if (tag[x] && tag[e[a].b])
                continue;
            if (pos[S] >= pos[e[a].b])
                continue;
            int y = e[a].b;
            res++;
            sum[1]++;
            sum[pos[S] + 1]--;
            sum[pos[y]]++;
            sum[tot + 1]--;
        } else if (!flag[e[a].b])
            dfs2(S, e[a].b);
    }
}
void dfs3(int S, int x) {
    if (check[x])
        return;
    flag[x] = 1;
    for (int a = la[x]; a; a = ff[a]) {
        if (check[e[a].b])
            continue;
        if (tag[e[a].b]) {
            if (tag[x] && tag[e[a].b])
                continue;
            if (pos[S] < pos[e[a].b])
                continue;
            int y = e[a].b;
            res++;
            sum[pos[y]]++;
            sum[pos[S] + 1]--;
        } else if (!flag[e[a].b])
            dfs3(S, e[a].b);
    }
}
int solve() {
    int po = 0;
    res = 0;
    tot = 0;
    memset(c, 0, sizeof(c));
    for (int i = 1; i <= n && !tot; i++)
        if (!c[i])
            dfs1(i);
    if (!tot)
        return 0;
    reverse(p + 1, p + tot + 1);
    memset(tag, 0, sizeof(tag));
    memset(pos, 0, sizeof(pos));
    memset(sum, 0, sizeof(sum));
    for (int i = 1; i <= tot; i++) {
        tag[p[i]] = 1;
        pos[p[i]] = i;
    }
    memset(flag, 0, sizeof(flag));

    for (int i = 1; i <= tot; i++)
        dfs2(p[i], p[i]);
    memset(flag, 0, sizeof(flag));
    for (int i = tot; i; i--)
        dfs3(p[i], p[i]);
    for (int i = 1; i <= tot; i++) {
        sum[i] += sum[i - 1];
        if (sum[i] == res) {
            po = p[i];
            break;
        }
    }
    check[po] = 1;
    return po;
}
int main() {
    int a, b;
    scanf("%d%d", &n, &m);
    for (int i = 1; i <= m; i++)
        scanf("%d%d", &a, &b), add(a, b);
    ans = solve();
    if (!ans) {
        printf("-1\n");
        return 0;
    }
    if (!solve())
        printf("%d\n", ans);
    else
        printf("-1\n");
    return 0;
}