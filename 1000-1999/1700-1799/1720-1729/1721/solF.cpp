#include<bits/stdc++.h>

using namespace std;

using LL=long long;

const int N=400010;

const int V=N;

const int E=N*3;

struct Edge

{

    int y,d,next;

}e[E<<1];

int h[V],last[V],cur[V],len=1,st,ed;

void ins(int x,int y,int d)

{

    int t=++len;

    e[t].y=y;e[t].d=d;e[t].next=last[x];last[x]=t;

}

void addEdge(int x,int y,int d)

{

    ins(x,y,d),ins(y,x,0);

}

int n1,n2,m,q,id[N],ex[N],ey[N];

bool bfs()

{

    memset(h,0,sizeof(h));h[st]=1;

    queue<int>q;

    q.push(st);

    cur[st]=last[st];

    while(!q.empty())

    {

        int x=q.front();q.pop();

        for(int i=last[x];i;i=e[i].next)

        {

            int y=e[i].y;

            if(e[i].d&&!h[y])

            {

                h[y]=h[x]+1;

                cur[y]=last[y];

                q.push(y);

                if (y==ed) return true;

            }

        }

    }

    return h[ed];

}

int dfs(int x,int f)

{

    if(x==ed)return f;

    int s=0,t;

    int &i=cur[x];

    for(;i;i=e[i].next)

    {

        int y=e[i].y;

        if(h[y]==h[x]+1&&e[i].d&&s<f)

        {

            t=dfs(y,min(f-s,e[i].d));

            s+=t;e[i^1].d+=t;e[i].d-=t;

            if (s==f) return s;

        }

    }

    if(!s)h[x]=0;

    return s;

}

int p[N],eid[N],u,v[N];

LL sum[N];

int main()

{

//    freopen("in","r",stdin);

//    freopen("out","w",stdout);

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    cin>>n1>>n2>>m>>q;

    st=n1+n2+1,ed=st+1;

    for(int i=1;i<=m;i++)

    {

        cin>>ex[i]>>ey[i];

        addEdge(ex[i],ey[i]+n1,1);

        id[i]=len;

    }

    for(int i=1;i<=n1;i++)addEdge(st,i,1);

    for(int i=1;i<=n2;i++)addEdge(i+n1,ed,1);

    int flow=0;

    while(bfs())flow+=dfs(st,INT_MAX);

    bfs();

    for(int i=1;i<=n1;i++)if(!h[i])p[++u]=i,v[i]=u;

    for(int i=1;i<=n2;i++)if(h[i+n1])p[++u]=i+n1,v[i+n1]=u;

    assert(u==flow);

    for(int i=1;i<=m;i++)

        if(e[id[i]].d)

        {

            if(v[ex[i]])eid[v[ex[i]]]=i;

            else eid[v[ey[i]+n1]]=i;

        }

    for(int i=1;i<=u;i++)sum[i]=sum[i-1]+eid[i];

    while(q--)

    {

        int o;

        cin>>o;

        if(o==1)

        {

            cout<<"1\n";

            if(p[u]<=n1)cout<<p[u]<<"\n";

            else cout<<-(p[u]-n1)<<"\n";

            u--;

            cout<<sum[u];

        }

        else

        {

            cout<<u<<"\n";

            for(int i=1;i<=u;i++)cout<<eid[i]<<" ";

        }

        cout<<endl;

    }

    return 0;

}