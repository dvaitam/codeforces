#include<set>

#include<cmath>

#include<queue>

#include<cstdio>

#include<cstring>

#include<iostream>

#include<algorithm>

#define inf (1<<30)

#define INF (1ll<<62)

#define prt(x) cout<<#x<<":"<<x<<" "

#define prtn(x) cout<<#x<<":"<<x<<endl

#define huh(x) printf("--------------case(%d)--------------\n",x)

#define st_ huh(1234)

#define en_ huh(4321)

#define ONLINE_JUDGE

using namespace std;

typedef long long ll;

typedef pair<int,int> ii;

///////////////////

template<class T>inline void Max(T &x,T y){if(x<y)x=y;};

template<class T>inline void Min(T &x,T y){if(x>y)x=y;};

template<class T>

void sc(T &x){

	x=0;char c;

	while(c=getchar(),c<48);

	do x=x*10+(c^48);

	while(c=getchar(),c>47);

}

template<class T>

void sim(T x){

	if(!x)return;

	sim(x/10);

	putchar('0'+x%10);

}

template<class T>

void pt(T x){

	if(!x)putchar('0');

	else sim(x);

	putchar('\n');

}

///////////////////

const int maxt=20005;

const int maxn=55;

const int maxm=105;

const int maxs=65536;

const double pi=acos(-1);

struct comp{

	double x,y;

	comp(){}

	comp(double x,double y):x(x),y(y){}

	comp operator+(const comp &a)const{

		return comp(x+a.x,y+a.y);

	}

	comp operator-(const comp &a)const{

		return comp(x-a.x,y-a.y);

	}

	comp operator*(const comp &a)const{

		return comp(x*a.x-y*a.y,y*a.x+x*a.y);

	}

}A[maxs],B[maxs],wm[maxs];

int N,rev[maxs];

void init(int len){

	int l=0;

	for(N=1;N<=len;N<<=1)l++;

	for(int i=1;i<N;i++)rev[i]=(rev[i>>1]>>1)|((i&1)<<l-1);

}

void dft(comp *a,int f){

	for(int i=0;i<N;i++)if(rev[i]<i)swap(a[i],a[rev[i]]);

	for(int m=1;m<N;m<<=1){

		wm[0]=comp(1,0);comp wn(cos(pi/m),f*sin(pi/m));

		for(int j=1;j<m;j++)wm[j]=wm[j-1]*wn;

		for(int j=0;j<N;j+=(m<<1)){

			for(int k=0;k<m;k++){

				comp x=a[j+k],y=a[j+k+m]*wm[k];

				a[j+k]=x+y;a[j+k+m]=x-y;

			}

		}

	}if(f==-1)for(int i=0;i<N;i++)a[i].x/=N;

}

int n,m,t,x;

int u[maxm],v[maxm],w[maxm];

int dis[maxn];

double p[maxm][maxt];

double psum[maxm][maxt];

double g[maxm][maxt];

double f[maxn][maxt];



void calc(int l,int r,int idx){

	int v=::v[idx];

	int mid=l+r>>1;

	int na=r-l,nb=mid-l;

	//p[][],f[][]

	int nn=na+nb;

	init(nn);

	for(int i=0;i<=na;i++)

		A[i]=comp(p[idx][i],0);

	for(int i=na+1;i<N;i++)

		A[i]=comp(0,0);

	//[0,r-l]

	for(int i=0;i<=nb;i++)

		B[i]=comp(f[v][i+l],0);

	for(int i=nb+1;i<N;i++)

		B[i]=comp(0,0);

	//[l,mid]

	dft(A,1);dft(B,1);

	for(int i=0;i<N;i++)A[i]=A[i]*B[i];

	dft(A,-1);

	//[l,r-l+mid]->[0,nn]

	//[mid,r-l+mid]->[nb,nn]

	for(int i=nb+1,j=mid+1;i<=nn&&j<=r;i++,j++)

		g[idx][j]+=A[i].x;

}

void divide(int l,int r){

	if(l==r){

		for(int i=0;i<m;i++){

			int u=::u[i],v=::v[i],w=::w[i];

			Min(f[u][l],g[i][l]+w+psum[i][l+1]*dis[v]);

		}

		return;

	}

	int mid=l+r>>1;

	divide(l,mid);

	for(int i=0;i<m;i++)calc(l,r,i);

	divide(mid+1,r);

}

int main(){

	#ifndef ONLINE_JUDGE

	freopen("data.in","r",stdin);

	freopen("data.out","w",stdout);

	#endif

	sc(n);sc(m);sc(t);sc(x);

	for(int i=0;i<m;i++){

		sc(u[i]);sc(v[i]);sc(w[i]);

		for(int j=1;j<=t;j++){

			sc(p[i][j]);

			p[i][j]/=100000.0;

		}

		psum[i][t+1]=0;

		for(int j=t;j>=1;j--)

			psum[i][j]=psum[i][j+1]+p[i][j];

	}

	memset(dis,-1,sizeof(dis));

	dis[n]=0;

	for(int i=1;i<=n;i++){

		bool update=0;

		for(int j=0;j<m;j++){

			int u=::u[j],v=::v[j],w=::w[j];

			if(dis[v]==-1)continue;

			if(dis[u]==-1||dis[u]>dis[v]+w){

				dis[u]=dis[v]+w;

				update=1;

			}

		}

		if(!update)break;

	}

	for(int i=1;i<=n;i++){

		dis[i]+=x;

		if(i!=n)for(int j=0;j<=t;j++)

			f[i][j]=dis[i];

	}

	divide(0,t);//[0,t]



	printf("%f\n",f[1][t]);

	return 0;

}