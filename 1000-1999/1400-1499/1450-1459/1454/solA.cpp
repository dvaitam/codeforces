#include <iostream>

#include <bits/stdc++.h>

#define pb push_back

#define yes cout << "YES\n"

#define no cout << "NO\n"

#define ll long long

 

using namespace std;

 

int main(int argc, char** argv) {

	ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

    	ll t;cin>>t;while(t--){

    	ll n;cin>>n;

    	vector<ll>a;

    	for(ll i=1;i<=n;i++){

    		a.pb(i);

		}

		sort(a.begin(),a.end());

		reverse(a.begin(),a.end());

		if(n==2){

			cout<<2<<" "<<1<<endl;

		}

		else{

			for(ll i=1;i<=n;i++){

	    		if(i==n/2){

	    			swap(a[i],a[i+1]);

				}

			}

			for(ll i=0;i<n;i++){

				cout<<" ";

	    		cout<<a[i];

			}

			

			cout<<endl;	

		}

		

	}

	return 0;	

	

}