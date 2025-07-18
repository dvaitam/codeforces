#include <bits/stdc++.h>
#define ll long long
#define pir pair<ll,ll>
#define mkp make_pair
#define fi first
#define se second
#define pb push_back
using namespace std;
const ll maxn=2e5+10, inf=1e17;
ll n,res[maxn],len,deg[maxn],rt;
vector<ll>to[maxn];
void dfs(ll u,ll fa,ll chs){
	ll d=0;
	for(ll v:to[u])
		if(v!=fa&&deg[v]>1) ++d;
	if((d-(u==rt))>1){
		puts("No"); exit(0);
	}
	if(u==rt){ ll x=0; res[len=1]=rt;
		for(ll v:to[u])
			if(deg[v]>1){
				x=v; dfs(v,u,0); break;
			}
		for(ll v:to[u])
			if(deg[v]==1) res[++len]=v;
		for(ll v:to[u])
			if(deg[v]>1&&x!=v) dfs(v,u,1);
		return;
	}
	if(chs){
		res[++len]=u;
		for(ll v:to[u])
			if(v!=fa&&deg[v]>1) dfs(v,u,0);
		for(ll v:to[u])
			if(v!=fa&&deg[v]==1) res[++len]=v;
	} else{
		for(ll v:to[u])
			if(v!=fa&&deg[v]==1) res[++len]=v;
		for(ll v:to[u])
			if(v!=fa&&deg[v]>1) dfs(v,u,1);
		res[++len]=u;
	}
}
int main(){// freopen("p.in","r",stdin);
	scanf("%lld",&n);
	for(ll i=1;i<n;i++){
		ll u,v; scanf("%lld%lld",&u,&v);
		to[u].pb(v), to[v].pb(u);
		++deg[u], ++deg[v];
	}
	if(n==2){
		puts("Yes"); puts("1 2"); return 0;
	} rt=0;
	for(ll i=1;i<=n;i++)
		if(deg[i]>1){
			rt=i; break;
		}
	dfs(rt,0,0);
	puts("Yes");
	for(ll i=1;i<=n;i++) printf("%lld ",res[i]);
	return 0;
}