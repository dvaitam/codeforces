#include<bits/stdc++.h>

#define endl "\n"

typedef long double ld;

typedef long long ll;

using namespace std;

vector<int> arr;

void print(int z, int n, int u, int v, int w, string s){

	for(int i=0;i<n;i++){

		if(arr[i]%3==0 && u) s[i]='1',u--;

		else if(arr[i]%3==1 && v)s[i]='1',v--;

		else if(w)s[i]='1',w--;

	}

	cout<<z<<endl<<s<<endl;

}

void code(){

	int n; cin>>n;

	arr.clear(); arr.resize(n);

	int rem[3]{}; string s(n,'0');

	for(int i=0;i<n;i++){ 

		cin>>arr[i];

		rem[arr[i]%3]++;

	}

	int x=0,y=0,z=0;

	if(rem[0]<=n/2){

		x=rem[0];

		if(x+rem[1]<=n/2)y=rem[1],z=n/2-rem[1]-x;

		else y=n/2-x;

		print(0,n,x,y,z,s);

	}

	else if(rem[0]==0||rem[0]==n){

		if(rem[0]==n)x=n/2;

		else if(rem[1]<=n/2)y=rem[1],z=n/2-rem[1];

		else y=n/2;

		print(1,n,x,y,z,s);

	}

	else if(rem[0]>=n/2){

		x=n/2;

		print(2,n,x,y,z,s);

	}

	else cout<<-1<<endl;

}



int main(){

	cin.tie(0)->sync_with_stdio(0);

	int t = 1; while(t--)code();

	return 0;

}