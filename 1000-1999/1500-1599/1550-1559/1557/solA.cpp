#include<bits/stdc++.h>
#define ll long long
#define ull unsigned long long
#define For(i,n,m) for(int i=n;i<=m;i++)
#define FOR(i,n,m) for(int i=n;i>=m;i--)
#define pb push_back
#define mp make_pair
#define fi first
#define se second
#define Mod 998244353
#define mo 1000000007
using namespace std;
void read(int &x){int ret=0;char c=getchar(),last=' ';while(!isdigit(c))last=c,c=getchar();while(isdigit(c))ret=ret*10+c-'0',c=getchar();x=last=='-'?-ret:ret;}
void readll(ll &x){ll ret=0;char c=getchar(),last=' ';while(!isdigit(c))last=c,c=getchar();while(isdigit(c))ret=ret*10+c-'0',c=getchar();x=last=='-'?-ret:ret;}
int pow(int x,int p,int mod){int ret=1;while(p){if(p&1)ret=1ll*ret*x%mod;x=1ll*x*x%mod;p>>=1;}return ret;}
int gcd(int a,int b){return b?gcd(b,a%b):a;}
ll gcdll(ll a,ll b){return b?gcdll(b,a%b):a;}

int t,n;
ll sum,maxn;

int main()
{
	//std::ios::sync_with_stdio(false),cin.tie(0);
	read(t);
	while(t--){
		read(n);
		int x;
		sum=0;
		maxn=-1e9-5;
		For(i,1,n){
			read(x);
			if(x>maxn)maxn=x;
			sum+=x;
		}
		sum-=maxn;
		printf("%.8lf\n",1.0*sum/(n-1)+maxn);
	}
	return 0;
}
/*
Is the original MAX(MIN) big(small) enough?
IOS IS USED!   Are you using read() and cin at the same time?
IOS IS USED!   Are you using printf and cout at the same time?

Are there multiple samples?
*/