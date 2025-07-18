#include<bits/stdc++.h>
#define Int int
#define pii pair<int,int>
#define fi first
#define se second
// #define int long long
#define mid (l+((r-l)>>1))
#define ls (p<<1)
#define rs (ls|1)
#define lt ls,l,mid
#define rt rs,mid+1,r
using namespace std;
const int N=2e6+5;
int b[N],n,cnt,a[N],ans[N];



inline int read(){
    int a=1,b=0;
    char ch=getchar();
    while(ch<'0'||ch>'9'){
        if(ch=='-') a=-a;
        ch=getchar();
    }
    while(ch>='0'&&ch<='9'){
        b=(b<<1)+(b<<3)+(ch^48);
        ch=getchar();
    }
    return a*b;
}

struct SegmentTree{
    pii d[N];
    inline void up(int p){
        d[p]=max(d[ls],d[rs]);
    }

    inline void update(int p,int l,int r,int x){
        if(l==r){
            d[p]=make_pair(b[x],x);
            return;
        }
        if(x<=mid) update(lt,x);
        else update(rt,x);
        up(p);
    }

    inline pii ask(int p,int l,int r,int x,int y){
        if(x>r||l>y) return make_pair(0,0);
        if(x<=l&&r<=y){
            return d[p];
        }
        return max(ask(lt,x,y),ask(rt,x,y));
    }

    
}t;

inline void topu(int u){
    int v=b[u];b[u]=0;
    t.update(1,1,n,u);//清空操作
    if(v!=n+1&&b[v]) topu(v);//表示先搜比我小的
    while(1){
        pii k=t.ask(1,1,n,1,a[u]-1);
        if(k.fi<=u) break;
        topu(k.se);
    }
    ans[u]=++cnt;
}

signed main(){
    n=read();
    for(int i=1;i<=n;i++){
        a[i]=read();
        if(a[i]!=-1) b[a[i]]=i;
    }    
    for(int i=1;i<=n;i++){
        if(a[i]==-1) a[i]=n+1;
        if(!b[i]) b[i]=n+1;
        t.update(1,1,n,i);
    }
    for(int i=n;i>=1;i--){
        if(!ans[i]) topu(i);
    }
    for(int i=1;i<=n;i++) printf("%d ",ans[i]);
}