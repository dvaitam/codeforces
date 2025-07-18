#pragma GCC optimize("Ofast")
#include <bits/stdc++.h>
using namespace std;
#define ll long long
#define pb push_back
#define mp make_pair
#define mt make_tuple
#define pii pair<int,int>
#define pll pair<ll,ll>
#define ldb double
template<typename T>void ckmn(T&a,T b){a=min(a,b);}
template<typename T>void ckmx(T&a,T b){a=max(a,b);}
void rd(int&x){scanf("%i",&x);}
void rd(ll&x){scanf("%lld",&x);}
void rd(char*x){scanf("%s",x);}
void rd(ldb&x){scanf("%lf",&x);}
void rd(string&x){scanf("%s",&x);}
template<typename T1,typename T2>void rd(pair<T1,T2>&x){rd(x.first);rd(x.second);}
template<typename T>void rd(vector<T>&x){for(T&i:x)rd(i);}
template<typename T,typename...A>void rd(T&x,A&...args){rd(x);rd(args...);}
template<typename T>void rd(){T x;rd(x);return x;}
int ri(){int x;rd(x);return x;}
template<typename T>vector<T> rv(int n){vector<T> x(n);rd(x);return x;}
template<typename T>void ra(T a[],int n,int st=1){for(int i=0;i<n;++i)rd(a[st+i]);}
template<typename T1,typename T2>void ra(T1 a[],T2 b[],int n,int st=1){for(int i=0;i<n;++i)rd(a[st+i]),rd(b[st+i]);}
template<typename T1,typename T2,typename T3>void ra(T1 a[],T2 b[],T3 c[],int n,int st=1){for(int i=0;i<n;++i)rd(a[st+i]),rd(b[st+i]),rd(c[st+i]);}
void re(vector<int> E[],int m,bool dir=0){for(int i=0,u,v;i<m;++i){rd(u,v);E[u].pb(v);if(!dir)E[v].pb(u);}}
template<typename T>void re(vector<pair<int,T>> E[],int m,bool dir=0){for(int i=0,u,v;i<m;++i){T w;rd(u,v,w);E[u].pb({v,w});if(!dir)E[v].pb({u,w});}}

const int N=150;
int a[N];

int cif(ll x){
	int cnt=0;
	while(x>0)x/=10,cnt++;
	return cnt;
}

void Brute(int n){
	int z=n/2;
	map<array<int,10>,vector<ll>> all;
	for(ll i=1;cif(i*i)<=n;i++){
		if(cif(i*i)==n){
			ll now=i*i;
			array<int,10> cnt;
			for(int j=0;j<10;j++)cnt[j]=0;
			while(now>0){
				cnt[now%10]++;
				now/=10;
			}
			all[cnt].pb(i*i);
		}
	}
	for(auto it:all){
		if(it.second.size()>=n){
			for(int i=0;i<n;i++)printf("%lldL%c",it.second[i],i==n-1?'}':',');
			return;
		}
	}
}
int main(){
	/*for(int i=1;i<=13;i+=2){
		printf("ans[%i]={",i);
		Brute(i);
		printf(";\n");
	}*/
	map<int,vector<ll>> ans;
	ans[1]={9L};
	ans[3]={169L,196L,961L};
	ans[5]={16384L,31684L,36481L,38416L,43681L};
	ans[7]={1493284L,3214849L,3912484L,4239481L,4293184L,4932841L,9132484L};
	ans[9]={236759769L,297769536L,369677529L,526977936L,677925369L,769729536L,773562969L,796763529L,927567936L};
	ans[11]={48458977956L,54785487969L,56487979584L,59988745476L,64755998784L,68597895744L,77956548849L,85745894976L,95887457649L,95896747584L,97598758464L};
	ans[13]={5898894567696L,6559678848969L,6586598874969L,6599894588676L,6788958546969L,6795885899664L,6965988947856L,8579966588964L,8585696978496L,8594788665969L,8698865979456L,8995896467856L,9579866858496L};

	for(int t=ri();t--;){
		int n=ri();
		if(n<=13){
			for(ll x:ans[n])printf("%lld\n",x);
		}else{
			int z=n/2;
			int found=0;
			for(int x=0;x<z;x++){
				for(int y=x+1;y<z;y++){
					if(x+z==y*2)continue;
					for(int i=0;i<n;i++)a[i]=0;
					a[x*2]+=1;
					a[y*2]+=1;
					a[z*2]+=1;
					a[x+y]+=2;
					a[x+z]+=2;
					a[y+z]+=2;
					for(int i=n-1;i>=0;i--)printf("%i",a[i]);
					printf("\n");
					found++;
					if(found==n)break;
				}
				if(found==n)break;
			}
		}
	}
	return 0;
}