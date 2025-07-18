// LUOGU_RID: 154446258
#include<bits/stdc++.h>
#include<iostream>
#include<iomanip>
#include<cstdio>
#include<cstring>
#include<vector>
#include<map>
#include<stack>
#include<math.h>
#include<stdlib.h>
#include<queue>
#include<set>
#include<string.h>
#include<string>
#include<stdio.h>
#include<deque>
using namespace std;
#define int long long
#define lll __uint128_t
#define ld long double
#define ull unsigned long long
#define fi first
#define sc second
#define all(v) v.begin(),v.end()
#define lowbit(x) (x&-x)
#define fixed(x) fixed<<setprecision(x)
#define io ios::sync_with_stdio(false),cin.tie(0),cout.tie(0)
//priority_queue<int,vector<int>,greater<int>>q;//从小到大
const int N = 2e5 + 10, base = 13331, mod = 1e9 + 7;
const double pie=3.1415926535897932,eps=1e-8;
int dir[]={0,1,0,-1,0};
const int dx[]={1,1,1,0,0,-1,-1,-1};
const int dy[]={1,0,-1,1,-1,0,1,-1};
typedef pair<int, int>PAIR;
void solve()
{
	int n;cin>>n;
	if(n==1){
		cout<<1<<endl;
		return ;
	}
	if(n%2){
		cout<<-1<<endl;
		return ;
	}
	cout<<n<<" ";
	for(int i=n-1,j=2;j<=n-2;i-=2,j+=2){
		cout<<i<<" "<<j<<" ";
	}
	cout<<1<<endl;
}
signed main()
{
    io;
    int T = 1;
    cin >> T;
    while(T -- ) solve();
}