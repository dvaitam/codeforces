#include<bits/stdc++.h>
using namespace std;
int main(){
	int t;
	cin>>t;
	while(t--){
		long long x,y,z,k,p,m=0;
		cin>>x>>y>>z>>k;
		for(long long i=1;i<=x;i++){
			for(long long j=1;j<=y;j++){
				p=k/(i*j);
				if(p*i*j==k&&p<=z){
					m=max(m,(x-i+1)*(y-j+1)*(z-p+1));
				}
			}
		}
		cout<<m<<"\n";
	}
}