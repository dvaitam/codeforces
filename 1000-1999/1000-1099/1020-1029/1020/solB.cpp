//Coded by Harsh

#include <bits/stdc++.h>
using namespace std;
#define f(i,n,x) for(int i=x; i<n; ++i)
#define fi(i,n,x) for(int i=x; i>n; --i)
#define f_it(it, v) for(it=v.begin(); it!=v.end(); ++it)
#define ARRAY_SIZE(arr) sizeof(arr)/sizeof(arr[0])
#define MOD 1000000007
#define ll long long int
#define F first
#define S second
#define mp make_pair
#define pb push_back
typedef pair<int,int> pii;
typedef vector<int> vi;

int n;
int p[10005];


int main(){

	cin>>n;

	f(i, n+1, 1){
		cin>>p[i];
	}

	f(i, n+1, 1){
		bool visited[n+1] = {false};

		int x = p[i];
		visited[i] = true;

		while(visited[x] != true){
			visited[x] = true;
			x = p[x];
		}

		cout<<x<<" ";
	}
	

	return 0;
}