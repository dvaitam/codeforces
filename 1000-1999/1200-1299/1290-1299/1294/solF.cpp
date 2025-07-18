#include <bits/stdc++.h>
using namespace std;
#define mem(a,b) memset(a,b,sizeof(a))
#define cin(a) scanf("%d",&a)
#define pii pair<int,int>
#define ll long long
#define gcd __gcd
const int inf = 0x3f3f3f3f;
const int maxn = 501110;
const int M = 1e9+7;

int head[maxn],to[maxn],Next[maxn],cnt = 2;

void add(int u,int v)
{
    to[cnt] = v;Next[cnt] = head[u];head[u] = cnt;cnt++;
}

int d[maxn],p,ans;

int fa[maxn];

void dfs(int u,int pre)
{   
    fa[u] = pre;
    d[u] = d[pre]+1;
    if(ans < d[u])
    {
        p = u;
        ans = d[u];
    }
    for(int i = head[u]; i ; i = Next[i])
    {
        int v = to[i];
        if(v == pre) continue;
        dfs(v,u);
    }
}

void find(int x)
{
    ans = 0;
    dfs(x,0);
}

struct node
{
    int a,b;
    node(int x,int y)
    {
        a = x;b = y;
    }
};

bool vis[maxn];

int main()
{
#ifdef ONLINE_JUDGE
#else
    freopen("data.in", "r", stdin);
#endif
    int n,a,b,c;
    cin(n);
    for(int i = 1,x,y; i < n; i++) 
    {
        scanf("%d%d",&x,&y);
        add(x,y);add(y,x);
    }
    find(1);
    a = p;
    ans = 0;
    find(a);
    b = p;
    queue<node> q;
    int x = b;
    while(x)
    {   
        q.push(node(x,0));
        vis[x] = 1;
        x = fa[x];
    }
    int u = 0,dis = 0;
    while(!q.empty())
    {
        u = q.front().a;
        dis = q.front().b;
        q.pop();
        for(int i = head[u]; i ; i = Next[i])
        {
            int v = to[i];
            if(vis[v]) continue;
            vis[v] = 1;
            q.push(node(v,dis+1));
        }
    }
    if(u == a || u == b)
    {
        for(int i = 1; i <= n; i++) 
        {
            if(i != a && i != b)
            {
                u = i;
                break;
            }
        }
    }
    cout<<ans-1+dis<<endl;
    cout<<a<<' '<<b<<' '<<u<<endl;
    return 0;
}