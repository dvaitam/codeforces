#include<bits/stdc++.h>

#define ll long long

#define pb push_back

#define inf INFINITY

const int N = 1e5;

using namespace std;

/*bool palindrome(string &s) {

    int n = s.length();

    for(int i = 0; i < n; i++) {

        if(s[i] != s[n - i - 1]) return false;

    }

    return true;

}*/

ll n,m,t,x,f;

int main(){

	ios_base::sync_with_stdio(0);

	cin.tie(0);cout.tie(0);

	cin>>t;

	while(t--){

	cin>>n>>m;

	ll a[n+1],b[m+1];

	for(int i=1;i<=n;i++){

		cin>>a[i];

	}

	for(int i=1;i<=m;i++){

		cin>>b[i];

	}

	sort(a+1,a+n+1);

	sort(b+1,b+m+1);

	for(int i=1;i<=n;i++){

		for(int j=1;j<=m;j++){

			if(a[i]==b[j]){

				f=1;

				x=a[i];

				break;

			}

		}

		if(f==1){

			break;

		}

	}

	if(f==1){

		cout<<"YES\n"<<1<<' '<<x<<"\n";

	}

	else{

		cout<<"NO\n";

	}

	f=0;

}

}

//(^w^)