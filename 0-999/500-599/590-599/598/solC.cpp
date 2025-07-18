#include<cstdio>

#include<cstring>

#include<algorithm>

#include<cctype>

#include<ctime>

#include<cstdlib>

#include<string>

#include<queue>

#include<cmath>

#define Rep(x,a,b) for (int x=a;x<=b;x++)

#define Per(x,a,b) for (int x=a;x>=b;x--)

#define ll long long

#define ld long double

using namespace std;

inline int IN(){

	int x=0,ch=getchar(),f=1;

	while (!isdigit(ch)&&(ch!='-')&&(ch!=EOF)) ch=getchar();

	if (ch=='-'){f=-1;ch=getchar();}

	while (isdigit(ch)){x=(x<<1)+(x<<3)+ch-'0';ch=getchar();}

	return x*f;

}

int n;

struct Vector{

	int idx;

	ld FUCK_CF;

	bool operator <(const Vector&x)const{return FUCK_CF<x.FUCK_CF;}

}a[100005];

int main(){

	n=IN();

	Rep(i,1,n){

		a[i].idx=i;

		int x=IN(),y=IN();

		a[i].FUCK_CF=atan2(1.0*y,1.0*x);

	}

	sort(a+1,a+n+1);

	int Ans1=a[1].idx,Ans2=a[n].idx;

	ld Min=a[1].FUCK_CF+2*acos(-1)-a[n].FUCK_CF;

	Rep(i,1,n-1){

		ld Ang=a[i+1].FUCK_CF-a[i].FUCK_CF;

		if (Ang<Min) Min=Ang,Ans1=a[i].idx,Ans2=a[i+1].idx;

	}

	printf("%d %d\n",Ans1,Ans2);

}