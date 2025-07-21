#include <bits/stdc++.h>
#define pb push_back
#define fst first
#define snd second
#define fore(i,a,b) for(ll i=a,jet=b;i<jet;++i)
#define SZ(x) ((int)x.size())
#define ALL(x) x.begin(),x.end()
#define mset(a,v) memset((a),(v),sizeof(a))
#define FIN ios::sync_with_stdio(0);cin.tie(0);cout.tie(0)
#define imp(v) {for(auto gdljh:v)cout<<gdljh<<" "; cout<<"\n";}
using namespace std;
typedef long long ll;
typedef pair<ll,ll> ii;
const ll MAXN=2e5+5;
ll a[MAXN];
typedef ll node;
node NEUT[3]={-1,0,-1};
node oper(node i, node j, ll ty){
	if(ty==1)return i+j;
	if(ty==0){
		if(i==-1)return j;
		if(j==-1)return i;
		return (a[i]>=a[j]?i:j);
	}
	return max(i,j);
}
struct STree{
	int n,ty; vector<node>t;
	STree(int n, int ty):n(n),ty(ty),t(2*n+5,NEUT[ty]){}
	void upd(int p, node v){
		for(p+=n,t[p]=v;p>1;p>>=1)t[p>>1]=oper(t[p],t[p^1],ty);
	}
	node query(int l, int r){
		node res=NEUT[ty];
		for(l+=n,r+=n;l<r;l>>=1,r>>=1){
			if(l&1)res=oper(res,t[l++],ty);
			if(r&1)res=oper(res,t[--r],ty);
		}
		return res;
	}
};
map<ll,ii>dp[MAXN];
int main(){FIN;
	ll t; cin>>t;
	while(t--){
		ll n,q; cin>>n>>q;
		fore(i,0,n)cin>>a[i];
		STree stm(n,0),sts(n,1),tim(n,2);
		fore(i,0,n)sts.upd(i,a[i]),stm.upd(i,i);
		
		
		ll cnt=0;
		auto f=[&](ll l, ll r, auto &&f)->ll{
			if(r-l<=0)return -1;
			ll vis=dp[l].count(r);
			auto &res=dp[l][r];
		
			if(vis&&tim.query(l,r)<=res.snd)return res.fst;
			res.snd=cnt;
			ll p=stm.query(l,r),s=sts.query(l,r);
			
			if(s-a[p]>a[p])return res.fst=r-l;
			return res.fst=max(f(l,p,f),f(p+1,r,f));
		};
		while(q--){
			ll ty,l,r; cin>>ty>>l>>r; l--;
			if(ty==1){
				cout<<f(l,r,f)<<"\n";
			}
			else {
				a[l]=r;
				sts.upd(l,r);
				stm.upd(l,l);
				tim.upd(l,cnt);
			}
			cnt++;
			
		}
		fore(i,0,n)dp[i].clear();
	}
	return 0;
}