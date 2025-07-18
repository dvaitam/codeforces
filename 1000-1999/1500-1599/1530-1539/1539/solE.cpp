#include<bits/stdc++.h>

#define Cn const

#define CI Cn int&

#define N 100000

#define Pr pair<int,int>

#define mp make_pair

using namespace std;

namespace FastIO

{

	#define FS 100000

	#define Tp template<typename Ty>

	#define Ts template<typename Ty,typename... Ar>

	#define tc() (FA==FB&&(FB=(FA=FI)+fread(FI,1,FS,stdin),FA==FB)?EOF:*FA++)

	#define pc(c) (FC==FE&&(clear(),0),*FC++=c)

	int OT;char oc,FI[FS],FO[FS],OS[FS],*FA=FI,*FB=FI,*FC=FO,*FE=FO+FS;

	void clear() {fwrite(FO,1,FC-FO,stdout),FC=FO;}struct IO_Cl {~IO_Cl() {clear();}}CL;

	Tp void read(Ty& x) {x=0;while(!isdigit(oc=tc()));while(x=(x<<3)+(x<<1)+(oc&15),isdigit(oc=tc()));}

	Ts void read(Ty& x,Ar&... y) {read(x),read(y...);}

	Tp void write(Ty x) {while(OS[++OT]=x%10+48,x/=10);while(OT) pc(OS[OT--]);pc(' ');}

}using namespace FastIO;

int n,m,w[N+5],f[N+5][2],k[N+5],xl[N+5],xr[N+5],yl[N+5],yr[N+5];set<pair<int,int> > A,B;

int main()

{

	int i;for(scanf("%d%d",&n,&m),A.insert(mp(0,0)),B.insert(mp(0,0)),i=1;i<=n;++i)

	{

		read(k[i],xl[i],xr[i],yl[i],yr[i]);

		(k[i]<xl[i]||k[i]>xr[i])&&(A.clear(),0),(k[i]<yl[i]||k[i]>yr[i])&&(B.clear(),0);

		while(!A.empty()&&A.begin()->first<yl[i]) A.erase(A.begin());while(!A.empty()&&(--A.end())->first>yr[i]) A.erase(--A.end());

		while(!B.empty()&&B.begin()->first<xl[i]) B.erase(B.begin());while(!B.empty()&&(--B.end())->first>xr[i]) B.erase(--B.end());

		if(A.empty()&&B.empty()) return puts("No"),0;

		f[i][0]=A.empty()?-1:A.begin()->second,f[i][1]=B.empty()?-1:B.begin()->second;

		~f[i][0]&&(B.insert(mp(k[i],i)),0),~f[i][1]&&(A.insert(mp(k[i],i)),0);

	}

	int o=!~f[n][0],x=f[n][o];for(i=n;i;--i) i==x&&(x=f[x][o^=1]),w[i]=o;

	for(puts("Yes"),i=1;i<=n;++i) write(w[i]);return 0;

}