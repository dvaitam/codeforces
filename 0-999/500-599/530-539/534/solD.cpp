#include <bits/stdc++.h>
using namespace std;
int const N = 200200;
struct Node {
    int x, v;
    friend bool operator < (Node a, Node b) {
        return a.v < b.v;
    }
} p[N];
int top[N];
int ans[N];

int main() {
    int n; scanf("%d", &n);
    for (int i = 1; i <= n; ++i) {
        scanf("%d", &p[i].v);
        p[i].x = i;
    }   
    p[0].v = -1;
    sort(p + 1, p + n + 1);
    memset(top, 0xff, sizeof top);
    for (int i = 1; i <= n; ++i) {
        if (p[i].v - p[i - 1].v > 1) {
            puts("Impossible");
            return 0;
        }
        top[p[i].v] = i;
    }
    int m = 0, now = 0;
    for (int i = 0; i < n; ++i) {
        if (top[now] == -1) {
            break;
        }
        ans[m++] = p[top[now]].x;
        int t = p[top[now]].v;
        --top[now];
        if (p[top[now]].v != t) top[now] = -1;
        ++now;
        if (now >= 3 && top[now - 3] != -1 && top[now - 2] != -1 && top[now - 1] != -1) {
            now -= 3;
        }
    }
    while (now >= 3) {
        now -= 3;
        if (top[now] == -1) continue;
        ans[m++] = p[top[now]].x;
        top[now] = -1;
        now += 1;
        while (top[now] != -1) {
            ans[m++] = p[top[now]].x;
            top[now] = -1;
            now += 1;
        }
    }
    if (m != n) {
        puts("Impossible");
        return 0;
    }
    puts("Possible");
    for (int i = 0; i < n; ++i) {
        if (i == n - 1) printf("%d\n", ans[i]);
        else printf("%d ", ans[i]);
    }
    return 0;
}