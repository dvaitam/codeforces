#include<bits/stdc++.h>

#include <vector>

using namespace std;



#define ll long long







int main(){

	ios_base::sync_with_stdio(0);

	cin.tie(0);

	ll t; cin >> t;

	while(t--){

		ll n,a,b;

		cin >> n >> a >> b;



		string s="";

		ll i=0;

		for(ll i=0; i<b; i++) s+='a'+i;

		while((ll)s.length() != n){

			i%=a;

			s+=s[i];

			i++;

		}

		cout << s << "\n";



	}

	return 0;

}