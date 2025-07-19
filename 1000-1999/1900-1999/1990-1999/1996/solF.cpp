#include<bits/stdc++.h>
#define int long long
using namespace std;
 
const int N=2e5+5;
 
int T,n,k,cnt,ans;
int a[N],b[N];
 
inline int calc(int x){
	cnt=ans=0;
	for(int i=1;i<=n;++i){
		if(a[i]<x) continue;
		int now=(a[i]-x)/b[i]+1;
		cnt+=now;
		ans+=(a[i]+a[i]-b[i]*(now-1))*now/2;
	}
	return cnt;
}
 
signed main(){
	ios::sync_with_stdio(0);cin.tie(0),cout.tie(0);
	cin>>T;
	while(T--){
		cin>>n>>k;
		for(int i=1;i<=n;++i) cin>>a[i];
		for(int i=1;i<=n;++i) cin>>b[i];
		int l=0,r=1e9,mid;
		while(l<r){
			mid=(l+r+1)>>1;
			if(calc(mid)>=k) l=mid;
			else r=mid-1;
		}
		calc(l);
		cout<<ans-(cnt-k)*l<<"\n";
	}
	return 0;
}