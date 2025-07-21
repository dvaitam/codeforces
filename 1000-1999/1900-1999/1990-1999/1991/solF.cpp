#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
const ll _=5e5+5;
ll n,q,x,y,a[_],b[_],c[_],d[_],i,j;
inline bool p(ll l,ll r){
	ll n=r-l+1,t,i,j;
	if(n<6)return 0;
	for(i=0;i<n;i++)b[i]=a[i+l];
	sort(b,b+n);
	for(i=2;i<n;i++){
		for(t=j=0;j<i-1;j++)if(b[j]+b[i-1]>b[i])break;
		if(j==i-1)continue;
		for(j++;j<n;j++)if(j!=i&&j!=i-1)c[t++]=b[j];
		for(j=2;j<t;j++)if(c[j-2]+c[j-1]>c[j])return 1;
	}
	return 0;
}
int main(){
	ios::sync_with_stdio(0);cin.tie(0);cout.tie(0);
	cin>>n>>q;
	for(i=1;i<=n;i++)cin>>a[i];
	for(i=j=1;i<=n;i++){
		while(j<=n&&!p(i,j))j++;
		d[i]=j;
	}
	while(q--){
		cin>>x>>y;
		cout<<(y<d[x]?"NO\n":"YES\n");
	}
}