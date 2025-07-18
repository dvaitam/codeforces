#include <bits/stdc++.h>
#include <cstdio>
#include <algorithm>
#include <cstring>

#define inf 1000000000

char s[100010], t[100010];
int n, d[30][30];

int main()
{
/*  freopen("in", "r", stdin);
    freopen("out", "w", stdout);*/
    fgets(s, sizeof(s), stdin); fgets(t, sizeof(t), stdin);
    if (strlen(s) != strlen(t))
    {
        printf("-1");
        return 0;
    }
    scanf("%d\n", &n);
    for (int i = 0; i < 26; i++)
        for (int j = 0; j < 26; j++)
            if (i == j)
                d[i][j] = 0;
            else d[i][j] = inf;
    for (int i = 0; i < n; i++)
    {
        char u, v; int l;
        scanf("%c %c%d\n", &u, &v, &l);
        u -= 'a'; v -= 'a';
        d[u][v] = std::min(l, d[u][v]);
    }
    for (int k = 0; k < 26; k++)
        for (int i=0; i<26; ++i)
            for (int j=0; j<26; ++j)
                if (d[i][k] < inf && d[k][j] < inf)
                    d[i][j] = std::min(d[i][j], d[i][k] + d[k][j]);
    int ans = 0;
    int m = strlen(s);
    for (int i = 0; i < m; i++)
    {
        s[i] -= 'a'; t[i] -= 'a';
        if (d[s[i]][t[i]] > d[t[i]][s[i]])
            std::swap(s[i], t[i]);
        int l = d[s[i]][t[i]];
        for (int j = 0; j < 26; j++)
            l = std::min(l, d[s[i]][j] + d[t[i]][j]);
        if (l == inf)
        {
            printf("-1");
            return 0;
        }
        if (l == d[s[i]][t[i]])
            s[i] = t[i];
        else
        {
            for (int j = 0; j < 26; j++)
                if (l == d[s[i]][j] + d[t[i]][j])
                {
                    s[i] = j;
                    break;
                }
        }
        ans += l;
    }
    printf("%d\n", ans);
    for (int i = 0; i < m; i++)
        printf("%c", s[i] + 'a');
}