#include<bits/stdc++.h>
#include<fstream>
#define int long long int
#define pb push_back
#define f first
#define s second
#define fs first.second
#define ff first.first
#define sf second.first
#define ss second.second
#define pf push_front 
#define inf 20000000000000009 
#define mod 1000000009 
#define bits __builtin_popcountll    
using namespace std;
const int A = 911382323;
const int B = 972663749;
//#pragma GCC optimize("-O2")
//#pragma GCC optimize("Ofast")
//#pragma GCC target("avx,avx2,fma")
//#pragma GCC optimize("O3,unroll-loops")
//#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")
int a[2005],b[2005],c[2005],d[2005],col,ans,e,mn,numm,num,mx,mmx,sz,ssum,cost,all,f,g,curr,speed,currx,dep,start,maxdep,tim,val,dist,pref,global,par,k,lvl,n,m,nm,imp,type,len,t,i,j,q,idx,h,sum,last,root,diff,pr;
string s,s1,s2,s3,alph = "abcdefghijklmnopqrstuvwxyz",str;
int px1,py1,px2,py2,px3,py3,x,y,z,tt,lq,rq;
int l,r,ll,rr,p,pp,qq,curr1,curr2,dp[2005][2005];
int gcd(int a,int b){ if(b == 0) return a; else return gcd(b,a%b); }
set<string> st;
signed main(){
    ios_base::sync_with_stdio(0),cin.tie(NULL),cout.tie(NULL);
    //ifstream cin ("INPUT.txt");
    //ofstream cout ("OUTPUT.txt");
   cin>>n;
   for(i=1;i<=2*n;i++){
   		cin>>d[i];
   		c[d[i]] = i;
   }
   vector<pair<int,int> > v; idx = 2*n;
   for(i=2*n;i>=1;i--){
   		if(c[i] > idx) continue;
   		v.pb({i,idx-c[i]+1});
   		b[i] = idx-c[i]+1;
   		idx = c[i]-1;
   } reverse(v.begin(),v.end()); idx = inf;
   //for(i=0;i<v.size();i++) cout<<v[i].f<<" "<<v[i].s<<endl;
   for(i=0;i<v.size();i++) dp[i][0] = 1;
   for(i=0;i<v.size();i++){
        for(j=1;j<=n;j++){
        	if(i == 0){
        		dp[i][v[i].s] = 1;
			}
			else{
        		if(dp[i-1][j] == 1) dp[i][j] = 1;
        		if(j-v[i].s >= 0) if(dp[i-1][j-v[i].s] == 1) dp[i][j] = 1;
			}
			if(dp[i][n] == 1) idx = i;
		}
   }
//   for(i=0;i<v.size();i++){
//   		for(j=1;j<=n;j++){
//   			cout<<dp[i][j]<<" ";	
//	    } cout<<endl;
//   }
   if(idx == inf) cout<<-1<<endl;
   else{  
   		x = idx; y = n; //cout<<x<<" "<<y<<endl;
   		vector<int> vv; ans = 0;
   		while(ans != n){
   			if(x == 0){
   				vv.pb(v[x].f);		
   				break;
			}
   			if(dp[x-1][y] == 1) x--;
   			else y -= v[x].s , vv.pb(v[x].f) , ans += v[x].s , x--;
		} reverse(vv.begin(),vv.end());
//		for(i=0;i<vv.size();i++) cout<<vv[i]<<" ";
//		cout<<endl;
		for(i=0;i<vv.size();i++){ 
			for(j=c[vv[i]];j<c[vv[i]]+b[vv[i]];j++){
				cout<<d[j]<<" "; a[j] = 1;
			} 
		} cout<<endl;
		for(i=1;i<=2*n;i++){
			if(a[i] == 0) cout<<d[i]<<" ";
		} cout<<endl;
   }
} 
//
// dont jump over problems
// think
// dont forget pragmas
// use sets wisely
// dont forget abt power of prefix
// dont forget about unordered map