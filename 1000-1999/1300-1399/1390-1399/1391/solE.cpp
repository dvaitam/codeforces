// LUOGU_RID: 102240164
#include <bits/stdc++.h>

#define mem(a,b) memset(a,b,sizeof(a))

#define fre(z) freopen(z".in","r",stdin),freopen(z".out","w",stdout)

using namespace std;

typedef long long ll;

typedef unsigned long long ull;

typedef pair<int,int> Pair;

const double eps=1e-8;

const int inf=2139062143;

#ifdef ONLINE_JUDGE

static char buf[1000000],*p1=buf,*p2=buf,obuf[1000000],*p3=obuf;

#define getchar() p1==p2&&(p2=(p1=buf)+fread(buf,1,1000000,stdin),p1==p2)?EOF:*p1++

#endif

inline void qread(){}template<class T1,class ...T2>

inline void qread(T1 &a,T2&...b){

    register T1 x=0;register bool f=false;char ch=getchar();

    while(ch<'0') f|=(ch=='-'),ch=getchar();

    while(ch>='0') x=(x<<3)+(x<<1)+(ch^48),ch=getchar();

    x=(f?-x:x);a=x;qread(b...);

}

inline void dread(){}template<class T1,class ...T2>

inline void dread(T1 &a,T2&...b){

    register double w=0;register ll x=0,base=1;

    register bool f=false;char ch=getchar();

    while(!isdigit(ch)) f|=(ch=='-'),ch=getchar();

    while(isdigit(ch)) x=(x<<3)+(x<<1)+(ch^48),ch=getchar();

    w=(f?-x:x);if(ch!='.') return a=w,dread(b...);x=0,ch=getchar();

    while(isdigit(ch)) x=(x<<3)+(x<<1)+(ch^48),base*=10,ch=getchar();

    register double tmp=(double)(x/(double)base);w=w+(double)(f?-tmp:tmp);a=w;dread(b...);

}

template<class T> T qmax(T x,T y){return x>y?x:y;}

template<class T,class ...Arg> T qmax(T x,T y,Arg ...arg){return qmax(x>y?x:y,arg...);}

template<class T> T qmin(T x,T y){return x<y?x:y;}

template<class T,class ...Arg> T qmin(T x,T y,Arg ...arg){return qmin(x<y?x:y,arg...);}

template<class T> T randint(T l,T r){static mt19937 eng(time(0));uniform_int_distribution<T>dis(l,r);return dis(eng);}

const int MAXN=5e5+7;

int T,n,m,dep[MAXN],fa[MAXN];vector<int>e[MAXN],vec[MAXN];bitset<MAXN>vis;

void dfs(int u){

    for(auto v:e[u]){

        if(v==fa[u]||vis[v]) continue;

        vis[v]=true;fa[v]=u;dep[v]=dep[u]+1;dfs(v);

    }

}int main(){

    qread(T);int i,j;

    while(T--){

        qread(n,m);vis.reset();

        for(i=1;i<=m;i++){

            int u,v;qread(u,v);

            e[u].emplace_back(v);e[v].emplace_back(u);

        }dep[1]=vis[1]=true;dfs(1);int p=1;for(i=1;i<=n;i++) if(dep[i]>dep[p]) p=i;

        if(dep[p]>=ceil(n*0.5)){

            puts("PATH");printf("%d\n",dep[p]);

            while(p) printf("%d ",p),p=fa[p];putchar(10);

        }else{

            puts("PAIRING");

            for(i=1;i<=n;i++) vec[dep[i]].emplace_back(i);int cnt=0;

            for(i=1;i<=dep[p];i++) cnt+=vec[i].size()/2;printf("%d\n",cnt);

            for(i=1;i<=dep[p];i++) for(j=0;j<vec[i].size()-1;j+=2) printf("%d %d\n",vec[i][j],vec[i][j+1]);

            for(i=1;i<=dep[p];i++) vec[i].clear();

        }for(i=1;i<=n;i++) e[i].clear();

    }

    #ifndef ONLINE_JUDGE

    cerr<<"Time: "<<clock()<<endl;

    system("pause > null");

    #endif

    return 0;

}