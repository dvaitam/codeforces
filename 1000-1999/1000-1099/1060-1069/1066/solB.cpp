#include<bits/stdc++.h>
#define FIO ios::sync_with_stdio(false)
#define N 100000
#define INF 0x3f3f3f3f
#define LL long long int
using namespace std;
int a[N];
int dp[N];
int main(){
	#ifndef ONLINE_JUDGE
	//	freopen("input.txt","r",stdin);
	#endif
	FIO;
	int n,r;
	cin>>n>>r;
	for(int i=1;i<=n;i++){ 
		cin>>a[i];
	}
	int cnt=0;
	int now=1;
	while(now<=n){
		int temp=min(now+r-1,n);
		while(a[temp]==0&&temp>0&&temp>=now-r+1){
			temp--;
		}
	//	cout<<now<<' '<<temp<<endl;
		if(temp<now-r+1||temp==0){
			cout<<"-1";
			return 0;
		}
		else now=temp+r;
		cnt++;
	}
	cout<<cnt<<endl;
	
}