#include <bits/stdc++.h>
#include <ext/numeric>
#include <hash_map>
using namespace std;
using namespace __gnu_cxx;

#define oo			1e9
#define OO			1e18
#define EPS			1e-7
#define f			first
#define s			second
#define pi 			acos(-1.0)
#define ll			long long
#define ld			long double
#define ull			unsigned long long
#define sz(x)		(int)x.size()
#define all(x)		x.begin(),x.end()
#define rall(x)		x.rbegin(),x.rend()
#define popCnt(x)   __builtin_popcount(x)

int n,m;

struct event{
	int type,x,val;
	event(int _type=0,int _x=0,int _val=0):type(_type),x(_x),val(_val){}
	bool operator<(const event& a)const{
		if(x!=a.x)return x<a.x;
		return type<a.type;
	}
};

vector<event> arr;
int prob[111];

int main() {
#ifndef ONLINE_JUDGE
	freopen("input.txt", "rt", stdin);
	//freopen("output.txt", "wt", stdout);
#endif
	scanf("%d%d",&n,&m);
	for(int i=0;i<n;i++){
		int a,h,l,r;
		scanf("%d%d%d%d",&a,&h,&l,&r);
		arr.push_back(event(-1,a-h,100-l));
		arr.push_back(event(1,a-1,100-l));
		arr.push_back(event(-1,a+1,100-r));
		arr.push_back(event(1,a+h,100-r));
	}
	for(int i=0;i<m;i++){
		int b,z;
		scanf("%d%d",&b,&z);
		arr.push_back(event(0,b,z));
	}
	sort(all(arr));
	double ans=0;
	for(int i=0;i<sz(arr);i++){
		if(arr[i].type==-1){
			prob[arr[i].val]++;
		}else if(arr[i].type==1){
			prob[arr[i].val]--;
		}else{
			if(prob[0])continue;
			double cur=arr[i].val;
			for(int j=1;j<100;j++)
				cur*=pow(1.0*j/100,prob[j]);
			ans+=cur;
		}
	}
	cout<<fixed<<setprecision(12)<<ans;
}