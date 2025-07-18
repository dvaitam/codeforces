#include <iostream>

#include <cstring>

#include <iomanip>

#include <cstdio>

#include <algorithm>

#include <cmath>

#include <cstdlib>

#include <vector>

#define INF 2147483647

#define ll long long

#define PI acos(-1)

#define MO 1000000007

//#define int long long

using namespace std;

	int nn,rr; 

	double n,r,jiao,jiao2,zj,h,d,p,db;

	double s;

inline ll read()

{

    ll x=0,f=1;char ch=getchar();

    while(ch<'0'||ch>'9'){if(ch=='-')f=-1;ch=getchar();}

    while(ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}

    return x*f;

}

inline void write(ll x){

    if(x==0){putchar('0');return;}if(x<0)putchar('-'),x=-x;

    ll len=0,buf[20];while(x)buf[len++]=x%10,x/=10;

    for(int i=len-1;i>=0;i--)putchar(buf[i]+'0');return;

}

int main()

{

	scanf("%d%d",&nn,&rr);

	n=nn,r=rr;

	jiao=PI/n;

//	printf("%.8lf\n",(double)180/n); 

//	printf("%.8lf\n%.8lf",(double)180/n/180*PI,jiao);

	jiao2=jiao/2;

	zj=PI/n;

	h=sin(zj)*r;

	d=cos(zj)*r;

	p=PI/2-zj-jiao2;

	db=tan(p)*h;

	s=(h*d-h*db)*n;

	printf("%.10lf",s);

	return 0;

}