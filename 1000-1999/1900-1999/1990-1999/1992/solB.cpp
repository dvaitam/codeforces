#include<bits/stdc++.h>
using namespace std;
int main(){
	int t;cin>>t;
	while(t--){
		int n,k,x,mx=0;
		cin>>n>>k;
		for(int i=1;i<=k;i++)
			cin>>x,mx=max(mx,x);
		cout<<(n-mx)*2-k+1<<'\n';
	}
}