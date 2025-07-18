#include<bits/stdc++.h>
using namespace std;
const int MAXN=200005;
int n,m,k,tot,g[MAXN],a,b,dis[MAXN],ans;
char s[MAXN];
vector<string>v;
struct edge
{
    int to,next,ki;
}e[MAXN<<1];
void add_edge(int from,int to,int ki)
{
    e[++tot].to=to;
    e[tot].ki=ki;
    e[tot].next=g[from];
    g[from]=tot;
    return;
}
void dfs(int x)
{
    if(x==n+1)
    {
        ++ans;
        v.push_back(string(s+1));
        return;
    }
    for(int i=g[x];i;i=e[i].next)
    {
        if(dis[e[i].to]+1==dis[x])
        {
            s[e[i].ki]='1';
            dfs(x+1);
            s[e[i].ki]='0';
            if(ans==k)break;
        }
    }
    return;
}
queue<int>q;
int main()
{
    scanf("%d %d %d",&n,&m,&k);
    for(int i=1;i<=m;++i)
    {
        scanf("%d %d",&a,&b);
        add_edge(a,b,i);
        add_edge(b,a,i);
    }
    memset(dis,-1,sizeof(dis));
    dis[1]=0;
    q.push(1);
    while(!q.empty())
    {
        int x=q.front();
        q.pop();
        for(int i=g[x];i;i=e[i].next)
        {
            if(dis[e[i].to]==-1)
            {
                dis[e[i].to]=dis[x]+1;
                q.push(e[i].to);
            }
        }
    }
    for(int i=1;i<=m;++i)
    {
        s[i]='0';
    }
    dfs(2);
    printf("%d\n",ans);
    for(int i=0;i<ans;++i)
    {
        cout<<v[i]<<endl;
    }
    return 0;
}