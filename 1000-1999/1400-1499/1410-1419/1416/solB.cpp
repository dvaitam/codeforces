//https://www.luogu.com.cn/problem/CF1416B

#include<bits/stdc++.h>

#define int long long

#define INF 0x3f3f3f3f

using namespace std;

const int N=1e4+5;

int n,T,a[N];

class node {

	public:

	int id,val;

}q1[N],q2[N];

//-----------------------------------

int read() {

	int x=0,f=1;

	char ch=getchar();

	while(ch<'0'||ch>'9') {if(ch=='-') f=-1;ch=getchar();}

	while(ch>='0'&&ch<='9') {x=(x<<1)+(x<<3)+(ch^48);ch=getchar();}

	return x*f;

}

//-----------------------------------

void write(int x) {

    if(x<0) putchar('-'),x=-x;

    if(x>9) write(x/10);

    putchar(x%10+'0');

}

//-----------------------------------

int Rand() {

	default_random_engine r(static_cast<unsigned int>(time(0)));

	uniform_int_distribution<>R(1,1);

	return R(r);

}

//-----------------------------------

bool cmp1(node aa,node bb) {

	return aa.val<bb.val; 

}

bool cmp2(node aa,node bb) {

	return aa.val>bb.val;

}

//------------------------------------

void Main() {

	int num=0;

	n=read();

	for(int i=1;i<=n;i++)

	{

		a[i]=read();

		num+=a[i];

	}

	if(num%n) {

		puts("-1");

		return;

	}

	int bj=0;

	num/=n;

	printf("%lld\n",3*(n-1));

	for(int i=2;i<=n;i++)

	{

		int p=a[i]%i;

		printf("%lld %lld %lld\n",1,i,(i-p)%i);

		a[1]-=(i-p)%i,a[i]+=(i-p)%i;

		printf("%lld %lld %lld\n",i,1,a[i]/i);

		a[1]+=a[i],a[i]=0;

	}

	for(int i=2;i<=n;i++) printf("%lld %lld %lld\n",1,i,num);

}

//-----------------------------------

signed main()

{

	T=read();

	while(T--) Main();

	return 0;

}