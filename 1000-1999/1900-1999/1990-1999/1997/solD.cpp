#include<bits/stdc++.h>
using namespace std;
const int N=2e5+5;
int t,n,a[N];
vector<int>e[N];
int S(int u){
	if(e[u].empty())
		return a[u];
	int r=1e9;
	for(auto v:e[u])
		r=min(r,S(v));
	return a[u]<r&&u!=1?(r+a[u])/2:r;
}
int main(){
	for(cin>>t;t;t--){
		cin>>n;
		for(int i=1;i<=n;i++)
			cin>>a[i],e[i].clear();
		for(int i=2,x;i<=n;i++)
			cin>>x,e[x].push_back(i);
		cout<<a[1]+S(1)<<"\n";
	}
	return 0;
}