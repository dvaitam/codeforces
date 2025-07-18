#pragma GCC optimize(3)
#include<bits/stdc++.h>
#define pa pair<ll,ll>
#define CLR(a,x) memset(a,x,sizeof(a))
using namespace std;
typedef long long ll;
const int maxn=3e5+10;

inline char gc(){
	return getchar();
	static const int maxs=1<<16;static char buf[maxs],*p1=buf,*p2=buf;
    return p1==p2&&(p2=(p1=buf)+fread(buf,1,maxs,stdin),p1==p2)?EOF:*p1++;
}
inline ll rd(){
    ll x=0;char c=gc();bool neg=0;
    while(c<'0'||c>'9'){if(c=='-') neg=1;c=gc();}
    while(c>='0'&&c<='9') x=(x<<1)+(x<<3)+c-'0',c=gc();
    return neg?(~x+1):x;
}

struct Edge{
	int a,b,l,ne;
}eg[maxn*2];
int egh[maxn],ect=1;
int fae[maxn],N,M,K;
ll dis[maxn];
bool flag[maxn],used[maxn*2];
priority_queue<pa,vector<pa>,greater<pa> > q;

inline void adeg(int a,int b,int c){
	eg[++ect].a=a,eg[ect].b=b,eg[ect].l=c,eg[ect].ne=egh[a],egh[a]=ect;
}

inline void dfs(int x,int f){
	for(int i=egh[x];i;i=eg[i].ne){
		int b=eg[i].b;
		if(eg[fae[b]].a!=x) continue;
		if(!K) return;
		K--;
		used[i]=1;dfs(b,x);
	}
}

inline void dijkstra(){
	CLR(dis,127);dis[1]=0;
	q.push(make_pair(0,1));
	while(!q.empty()){
		int p=q.top().second;q.pop();
		if(flag[p]) continue;
		flag[p]=1;
		for(int i=egh[p];i;i=eg[i].ne){
			int b=eg[i].b;
			if(dis[b]>dis[p]+eg[i].l){
				dis[b]=dis[p]+eg[i].l;
				q.push(make_pair(dis[b],b));
				fae[b]=i;
			}
		}
	}
}

int main(){
    //freopen("","r",stdin);
    int i,j,k;
    N=rd(),M=rd(),K=rd();
    for(i=1;i<=M;i++){
    	int a=rd(),b=rd(),c=rd();
    	adeg(a,b,c);adeg(b,a,c);
    }
    dijkstra();
    printf("%d\n",min(N-1,K));
    dfs(1,0);
    for(i=1;i<=ect;i++){
    	if(used[i]) printf("%d ",i>>1);
    }
    return 0;
}