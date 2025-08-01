#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <algorithm>
#include <cmath>
using namespace std;
const int maxn=30;
int nf,ne,ns,df,de;
double rf,re,rs;

namespace Ninit{
	void init(){
		cin>>nf>>ne>>ns>>rf>>re>>rs>>df>>de;
		rf=sqrt(rf*rf-1),re=sqrt(re*re-1),rs=sqrt(rs*rs-1);
	}
}

namespace Nsolve{
	double ans;
	double Len(double a,double b,double c,double d){
		return max(min(b,d)-max(a,c),0.);
	}
	bool U[maxn];int a[maxn],m;double b[maxn];
	double calc(){
		int i,j;double Df,De;
		double Fc=2*nf*rf*df+2*ne*re*de;
		for(m=0,i=0;i<nf+ne+ns;++i)if(!U[i]){
			for(Df=De=0,j=0;j<ns;++j){
				Df+=df*Len(i/2-rf,i/2+rf,a[j]-rs,a[j]+rs);
				De+=de*Len(i/2-re,i/2+re,a[j]-rs,a[j]+rs);
			}
			Fc+=Df,b[m++]=De-Df;
		}
		sort(b,b+m),reverse(b,b+m);
		for(i=0;i<ne;++i)Fc+=b[i];
		return Fc;
	}
	void dfs(int x,int y){
		if(nf+ne+y<x)return;
		if(x==nf+ne+ns)
			return ans=max(ans,calc()),void();
		U[x]=0,dfs(x+1,y);
		if(y<ns && (~x&1 || U[x-1])){
			U[x]=1,a[y]=x/2;
			dfs(x+1,y+1);
		}
	}
	void solve(){
		dfs(0,0);
		printf("%.10lf\n",ans);
	}
}

int main(){
	//freopen("A.in","r",stdin);
	Ninit::init();
	Nsolve::solve();
	return 0;
}