#include <bits/stdc++.h>
#include <cstdio>
#include <cstring>
#include <algorithm>
using namespace std;

typedef long long LL;
const LL Inf = 100000000000LL;
const int MaxN = 300005;
int ax[MaxN],c[MaxN],d[MaxN],l[MaxN],n,m,K;
int st[MaxN], mx[MaxN][2], pre[MaxN][2];
LL sel[MaxN][2],ans[MaxN];
const LL B = 200000;
LL dp[MaxN][2];
int r[MaxN],mp[MaxN],M,N;
bool cmpR(int p1, int p2)
{
    return ax[p1]<ax[p2];
}
void Lisan()
{
    N = n*2+m;
    for(int i = 0; i < N; i++)r[i]=i;
    sort(r,r+N,cmpR);
    mp[1] = ax[r[0]];
    ax[r[0]] = M = 1;
    for(int i = 1; i < N; i++)
    {
        if(ax[r[i]]!=mp[M]) mp[++M] = ax[r[i]];
        ax[r[i]] = M;
    }
}

int main()
{
    scanf("%d",&n);
    for(int i = 0; i < n; i++)
    {
        scanf("%d%d",&c[i],&ax[i]);
        ax[i+n] = ax[i]-1;
    }
    scanf("%d",&m);
    for(int i = 0; i < m; i++)
    {
        scanf("%d%d",&d[i],&l[i]);
        ax[i+n*2] = l[i];
    }
    Lisan();
    memset(st,-1,sizeof(st));
    for(int i = 0; i < n; i++)
        st[ax[i+n]] = i;
    memset(mx,-1,sizeof(mx));
    for(int i = 0; i < m; i++)
    {
        int j = ax[i+2*n];
        if(mx[j][0]==-1 || d[i] > d[mx[j][0]])
        {
            mx[j][1] = mx[j][0];
            mx[j][0] = i;
        }
        else if(mx[j][1]==-1 || d[i] > d[mx[j][1]])
            mx[j][1] = i;
    }
    for(int i = 0; i <= M+1; i++)
        for(int j = 0; j < 2; j++)
            dp[i][j] = -Inf;
    dp[1][0] = 0;
    memset(sel,-1,sizeof(sel));
    for(int i = 1; i <= M; i++)
        for(int j = 0; j < 2; j++)if(dp[i][j]!=-Inf)
        {
            if(dp[i+1][0] < dp[i][j])
                dp[i+1][0] = dp[i][j], pre[i+1][0] = j, sel[i+1][0]=-1;
            if(j==0)
            {
                if(mx[i][0]!=-1 && d[mx[i][0]]>=c[st[i]])
                {
                    if(dp[i+1][0] < dp[i][j]+c[st[i]])
                    {
                        dp[i+1][0] = dp[i][j]+c[st[i]];
                        pre[i+1][0] = j;
                        sel[i+1][0] = st[i]*B+mx[i][0];
                    }
                }
            }
            else
            {
                int xx=-1;
                if(mx[i][1]!=-1 && d[mx[i][1]]>=c[st[i-1]])
                    xx = mx[i][0];
                else xx = mx[i][1];
                if(xx!=-1 && d[xx]>=c[st[i]])
                {
                    if(dp[i+1][0] < dp[i][j]+c[st[i]])
                    {
                        dp[i+1][0] = dp[i][j]+c[st[i]];
                        pre[i+1][0] = j;
                        sel[i+1][0] = st[i]*B+xx;
                    }
                }

            }
            if(mx[i+1][0]!=-1 && d[mx[i+1][0]]>=c[st[i]])
            {
                if(dp[i+1][1] < dp[i][j]+c[st[i]])
                {
                    dp[i+1][1] = dp[i][j]+c[st[i]];
                    pre[i+1][1] = j;
                    if(mx[i+1][1]!=-1 && d[mx[i+1][1]]>=c[st[i]])
                        sel[i+1][1] = st[i]*B+mx[i+1][1];
                    else sel[i+1][1] = st[i]*B+mx[i+1][0];
                }
            }
        }
    printf("%lld\n",dp[M+1][0]);
    K = 0;
    for(int i = M+1, j = 0; i >= 1; i--)
    {
        if(sel[i][j]!=-1)ans[K++] = sel[i][j];
        j = pre[i][j];
    }
    printf("%d\n",K);
    for(int i = 0; i < K; i++)
        printf("%lld %lld\n",ans[i]%B+1,ans[i]/B+1);

    return 0;
}