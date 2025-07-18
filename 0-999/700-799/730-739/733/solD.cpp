#include<set>

#include<map>

#include<queue>

#include<cmath>

#include<string>

#include<cstdio>

#include<vector>

#include<cstring>

#include<iostream>

#include<algorithm>

#define rep(i,a,b) for (int i=a; i<=b; i++)

using namespace std;

typedef long long LL;



inline int read() {

    int x=0,f=1; char ch=getchar();

    while (!(ch>='0'&&ch<='9')) {if (ch=='-')f=-1;ch=getchar();}

    while (ch>='0'&&ch<='9') {x=x*10+(ch-'0'); ch=getchar();}

    return x*f;

}



const int N = 100005;



struct rectangle {

	int a, b, c, id;

} rec[N];



bool cmp(rectangle x, rectangle y) {

	if (x.a==y.a&&x.b==y.b) return x.c>y.c;

	if (x.a==y.a) return x.b>y.b;

	return x.a>y.a;

}



int n, p1, p2;

double ans = 0;

bool flag=false;



int main() {



	#ifndef ONLINE_JUDGE

		freopen("data.in","r",stdin);

		freopen("data.out","w",stdout);

	#endif

	

	n=read(); rep(i,1,n) {

		rec[i].a=read(); rec[i].b=read(); rec[i].c=read(); rec[i].id=i;

		if (rec[i].a<rec[i].b) swap(rec[i].a,rec[i].b);

		if (rec[i].a<rec[i].c) swap(rec[i].a,rec[i].c);

		if (rec[i].b<rec[i].c) swap(rec[i].b,rec[i].c);

		if ((double)rec[i].c/2>ans) ans=(double)rec[i].c/2,p1=i;

	}

	sort(rec+1,rec+n+1,cmp);

	

	rep(i,2,n) {

		if (rec[i].a==rec[i-1].a&&rec[i].b==rec[i-1].b){

			if ((double)(min(rec[i].b,rec[i].c+rec[i-1].c))/2>ans) flag=1,ans=(double)(min(rec[i].b,rec[i].c+rec[i-1].c))/2,p1=rec[i].id,p2=rec[i-1].id;

		}

	}

	if (!flag) cout<<1<<endl<<p1<<endl; else cout<<2<<endl<<p1<<" "<<p2<<endl;

	

	return 0;

}