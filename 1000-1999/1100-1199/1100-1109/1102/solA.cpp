#include<bits/stdc++.h>
using namespace std;
#define ll long long
int main(){
	ll t;
	t=1;
//cin>>t;
	while(t--){
		ll n;
		cin>>n;
		if(n==0)
			cout<<0;
		else if(n==1)
			cout<<1;
		else if(n==2)
			cout<<1;
		else if(n==3)
			cout<<0;
		else if(n==4)
			cout<<0;
		else{
			ll sum1=0,sum2=0;
			ll a=n;
			if(a%4==0){
				cout<<0;
			}
			else if(a%2==0){
				cout<<1;
			}
			else{
				a=a-1;
				if(a%4==0){
					cout<<1;
				}
				else if(a%2==0){
					cout<<0;
				}
			}
		}
	}
		
}