#include<bits/stdc++.h>
using namespace std;
#define ll long long

int main(){
	ll n,t,a,b;
	cin>>t;
	while(t--){
		cin>>n>>a>>b;
		if(2*a<b) cout<<n*a<<"\n";
		else{
			cout<<(n/2)*b+(n%2)*a<<"\n";
		}
	}

	return 0;
}