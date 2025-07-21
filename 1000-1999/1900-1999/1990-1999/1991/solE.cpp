#include<bits/stdc++.h>
#define eb emplace_back
using namespace std;
typedef long long ll;
const ll _=1e5+5;
ll N,n,m,x,y,t,a[_],b[_],c[_],i;vector<ll>E[_],A,B;
inline void p1(ll x){
	for(auto y:E[x]){
		if(!a[y])a[y]=a[x]^1,p1(y);
		else if(a[y]==a[x])t=1;
	}
}
inline void P(){
	cin>>n>>m;A.clear();B.clear();
	for(i=0;i<=n;i++)a[i]=b[i]=c[i]=0,E[i].clear();
	while(m--)cin>>x>>y,E[x].eb(y),E[y].eb(x);
	a[1]=2;t=0;p1(1);
	if(t){
		cout<<"Alice"<<endl;
		for(i=0;i<n;i++){
			cout<<"1 2"<<endl;
			cin>>x>>y;
		}
	}
	else{
		cout<<"Bob"<<endl;
		for(i=1;i<=n;i++)(a[i]==2?A:B).eb(i);
		for(i=0;i<n;i++){
			cin>>x>>y;x^=y;
			if(x==1){
				if(A.size())cout<<A.back()<<" 2"<<endl,A.pop_back();
				else cout<<B.back()<<" 3"<<endl,B.pop_back();
			}
			else{
				if(B.size())cout<<B.back()<<" 1"<<endl,B.pop_back();
				else cout<<A.back()<<' '<<(x^1)<<endl,A.pop_back();
			}
		}
	}
}
int main(){cin>>N;while(N--)P();}