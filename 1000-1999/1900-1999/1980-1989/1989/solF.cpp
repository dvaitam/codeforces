#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
const int N=4e5+5;
int n,m,q,fa[N],sz[N];
ll ans;
struct edge{
        int x,y,t;
};
vector<int> e[N];
int dfn[N],low[N],st[N],tp,co[N],tot;
bool vis[N];
int gf(int x){
        return fa[x]==x?fa[x]:fa[x]=gf(fa[x]);
}
void merge(int x,int y){
        x=gf(x),y=gf(y);
        if(x==y)return;
        if(sz[x]>1)ans-=1ll*sz[x]*sz[x];
        if(sz[y]>1)ans-=1ll*sz[y]*sz[y];
        fa[y]=x;
        sz[x]+=sz[y];
        ans+=1ll*sz[x]*sz[x];
}
void tarjan(int x){
        dfn[x]=low[x]=++tot;
        st[++tp]=x;
        vis[x]=true;
        for(auto v:e[x]){
                if(!dfn[v]){
                        tarjan(v);
                        low[x]=min(low[x],low[v]);
                }else if(vis[v]){
                        low[x]=min(low[x],dfn[v]);
                }
        }
        if(low[x]==dfn[x]){
                co[x]=x;
                while(st[tp]!=x){
                        co[st[tp]]=x;
                        vis[st[tp]]=false;
                        --tp;
                }
                --tp;
                vis[x]=false;
        }
}
void solve(int l,int r,vector<edge>&ed){
        if(l==r){
                if(l>q) return;
                for(auto v:ed) merge(v.x,v.y);
                printf("%lld\n",ans);
                return;
        }
        int mid=l+r>>1;
        tot=0;
        vector<edge> el,er;
        for(auto &v:ed){
                v.x=gf(v.x);
                v.y=gf(v.y);
                e[v.x].clear();
                e[v.y].clear();
                dfn[v.x]=0;
                dfn[v.y]=0;
        }
        for(auto v:ed) if(v.t<=mid) e[v.x].push_back(v.y);
        for(auto v:ed){
                if(v.t<=mid){
                        if(!dfn[v.x]) tarjan(v.x);
                        if(!dfn[v.y]) tarjan(v.y);
                        if(co[v.x]==co[v.y]){
                                el.push_back(v);
                                continue;
                        }
                }
                er.push_back(v);
        }
        solve(l,mid,el);
        solve(mid+1,r,er);
}
int main(){
        scanf("%d%d%d",&n,&m,&q);
        for(int i=1;i<=n+m;++i) fa[i]=i, sz[i]=1;
        vector<edge> E;
        for(int i=1,x,y;i<=q;++i){
                scanf("%d%d",&x,&y);
                char c=getchar();
                while(c!='R' && c!='B') c=getchar();
                if(c=='R') E.push_back(edge{y+n,x,i});
                else E.push_back(edge{x,y+n,i});
        }
        solve(1,q+1,E);
        return 0;
}
