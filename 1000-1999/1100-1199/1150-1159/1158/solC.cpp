#include <cstdio>
#include <cstring>
#include <string>
#include <cmath>
#include <iostream>
#include <algorithm>
#include <map>
#include <set>
#include <queue>
#include <vector>
using namespace std;
typedef long long ll;
typedef unsigned int uint;
typedef unsigned long long ull;
typedef pair<int, int> PII;
typedef vector<int> VI;
#define fi first
#define se second
#define MP make_pair

int read()
{
    int v = 0, f = 1;
    char c = getchar();
    while (c < 48 || 57 < c) {if (c == '-') f = -1; c = getchar();}
    while (48 <= c && c <= 57) v = (v << 3) + v + v + c - 48, c = getchar();
    return v * f;
}

const int N = 501000;
VI V[N];
int ans[N], a[N], stk[N], n, cnt, init[N], Q[N], vis[N];

int main()
{
    int cas = read();
    while (cas--)
    {
        n = read();
        for (int i = 1; i <= n; i++)
        {
            a[i] = read();
            if (a[i] == -1)
                a[i] = i + 1;
        }
        int top = 0;
        stk[0] = n + 1;
        for (int i = 1; i <= n + 1; i++)
            vis[i] = 0;
        vis[n + 1] = 1;
        int flg = 0;
        for (int i = 1; i <= n + 1; i++)
            init[i] = 0;
        for (int i = n; i >= 1; i--)
        {
            if (!vis[a[i]])
            {
                flg = 1;
                break;
            }
            while (stk[top] != a[i])
            {
                vis[stk[top]] = 0;
                V[i].push_back(stk[top]);
                init[stk[top]]++;
                top--;
            }
            V[stk[top]].push_back(i);
            init[i]++;
            stk[++top] = i;
            vis[i] = 1;
        }
        if (flg)
            puts("-1");
        else
        {
            int H = 1, T = 0;
            for (int i = 1; i <= n + 1; i++)
                if (init[i] == 0)
                    Q[++T] = i;
            while (H <= T)
            {
                int u = Q[H++];
                for (int e = 0; e < V[u].size(); e++)
                {
                    int v = V[u][e];
                    init[v]--;
                    if (init[v] == 0)
                        Q[++T] = v;
                }
            }
            if (T != n + 1)
                puts("-1");
            else
            {
                for (int i = 2; i <= n + 1; i++)
                    ans[Q[i]] = n - i + 2;
                for (int i = 1; i <= n; i++)
                    printf("%d ", ans[i]);
                puts("");
            }
        }
        for (int i = 1; i <= n + 1; i++)
            VI().swap(V[i]);
    }
}