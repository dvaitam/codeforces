#include<bits/stdc++.h>
using namespace std;
const int N=5005,K=505;
typedef long long ll;
int T,n,r,k,a[N],dfn[N],out[N],rv[N],fa[N];
vector<int>e[N];
bool op[N],ban[N];
struct frac{
	ll x,y;
	frac operator +(frac a){return {x+a.x,y+a.y};}
	bool operator <(frac a){return x*a.y<y*a.x;}
	bool operator ==(frac a){return x*a.y==y*a.x;}
}s[N],w[N];
void dfs(int x,int y){
	dfn[x]=++dfn[0],rv[dfn[0]]=x,fa[x]=y;
	s[x]={a[x],1};
	for(auto v:e[x])if(v!=y)dfs(v,x),s[x]=s[x]+s[v];
	out[x]=dfn[0];
}
void calc(int x){
	int f1=0,f2;
	for(int i=x;i;i=fa[i])if(op[i])f1=i;
	if(!f1){
		w[x]={a[x],1};
		if(k){
			frac mn=w[x];
			for(int i=x;i;i=fa[i])if(!ban[i]&&(s[i]<mn||mn<w[x]&&s[i]==mn))f1=i,mn=s[i];
		}
		if(f1)--k,w[x]=s[f1],op[f1]=true;
		return;
	}
	vector<int>tmp;
	int rs=k+1;
	for(int i=dfn[f1];i<=out[f1];++i){
		int v=rv[i];
		if(v<x&&w[v]<frac{a[v],1}){
			f2=0;
			for(int o=v;out[o]<dfn[x]||dfn[o]>dfn[x];o=fa[o])if(!ban[o])f2=o;
			if(!f2){
				rs=-1;
				break;
			}
			tmp.emplace_back(f2),--rs;
			i=out[f2];
		}
	}
	if(rs<0){
		w[x]=s[f1];
		return;
	}
	w[x]={a[x],1},f2=0;
	frac mn=w[x];
	if(rs){
		for(int i=x;i;i=fa[i])if(!ban[i]&&(s[i]<mn||s[i]<w[x]&&s[i]==mn))f2=i,mn=s[i];
	}
	if(mn<s[f1]){
		op[f1]=false;
		for(auto v:tmp)op[v]=true;
		k=rs,w[x]=mn;
		if(f2)--k,op[f2]=true;
	}else w[x]=s[f1];
}
void solve(){
	scanf("%d%d%d",&n,&r,&k);
	for(int i=1;i<=n;++i)scanf("%d",&a[i]),e[i].clear(),op[i]=ban[i]=0;
	dfn[0]=0;
	for(int i=2,x,y;i<=n;++i){
		scanf("%d%d",&x,&y);
		e[x].emplace_back(y);
		e[y].emplace_back(x);
	}
	dfs(r,0);
	for(int i=1;i<=n;++i){
		calc(i);
		for(int j=i;j;j=fa[j])if(w[i]<s[j])ban[j]=true;
	}
	vector<int>tmp;
	for(int i=1;i<=n;++i)if(op[i])tmp.emplace_back(i);
	printf("%d\n",tmp.size());
	for(auto v:tmp)printf("%d ",v);
	putchar(10);
}
int main(){
	scanf("%d",&T);
	while(T--)solve();
	return 0;
}