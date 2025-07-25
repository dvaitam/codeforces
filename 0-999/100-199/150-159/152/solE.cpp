#include <iostream>
#include <cstdio>
#include <cstring>
#include <queue>
#include <algorithm>
using namespace std;
#define INF 1061109567
const int N=205;
const int M=300;
int dx[]={0,1,0,-1},dy[]={1,0,-1,0};

int n,m,K;
int dp[N][M],pre[N][M];
int mat[N],vis[N],st[N];
bool in[N][M];
queue<int> que;

void spfa()
{
    while(que.size())
    {
        int top=que.front(); que.pop();
        int x=top/1000/m,y=(top/1000)%m,s=top%1000;
        in[x*m+y][s]=0;
        for(int i=0;i<4;i++)
        {
            int tx=x+dx[i],ty=y+dy[i];
            if(tx>=n||tx<0||ty>=m||ty<0) continue;
            int ts=s|st[tx*m+ty];
            if(dp[tx*m+ty][ts]>dp[x*m+y][s]+mat[tx*m+ty])
            {

                dp[tx*m+ty][ts]=dp[x*m+y][s]+mat[tx*m+ty];
                pre[tx*m+ty][ts]=top;
                if(in[tx*m+ty][ts]==0&&s==ts)
                {
                    in[tx*m+ty][ts]=1;
                    que.push( (tx*m+ty)*1000+ts );
                }
            }
        }
    }
}
void getans(int x,int y,int mask)
{
    //cout<<x<<" "<<y<<" "<<mask<<endl;
    vis[x*m+y]=1;
    int tmp=pre[x*m+y][mask];
    if(tmp==0) return;
    int tx=tmp/1000/m,ty=(tmp/1000)%m,s1=tmp%1000;
    getans(tx,ty,s1);
    if(tx==x&&ty==y)
        getans(tx,ty,((mask-s1)|st[x*m+y]));
}
int main()
{
    scanf("%d%d%d",&n,&m,&K);
    for(int i=0;i<n;i++)
        for(int j=0;j<m;j++) scanf("%d",&mat[i*m+j]);

    memset(dp,63,sizeof(dp));
    int a,b;
    for(int i=0;i<K;i++)
    {
        scanf("%d%d",&a,&b);
        a--;b--;
        st[a*m+b]=(1<<i);
        dp[a*m+b][(1<<i)]=mat[a*m+b];
    }

    int mask=(1<<K)-1;
    for(int s=1;s<=mask;s++)
    {
        for(int i=0;i<n*m;i++)
        {
            if(st[i]&&!(st[i]&s)) continue;
            for(int p=(s-1)&s;p;p=(p-1)&s)
            {
                int s1=p|st[i],s2=(s-p)|st[i];
                int d=dp[i][s1]+dp[i][s2]-mat[i];
                if(d<dp[i][s])
                {
                    pre[i][s]=i*1000+s1;
                    dp[i][s]=d;
                }
            }
            if(dp[i][s]!=INF)
                in[i][s]=1,que.push(i*1000+s);
        }
        spfa();
    }
    printf("%d\n",dp[a*m+b][mask]);
    getans(a,b,mask);
    for(int i=0;i<n;i++,puts(""))
        for(int j=0;j<m;j++)
        {
            if(vis[i*m+j]) putchar('X');
            else putchar('.');
        }
    return 0;
}