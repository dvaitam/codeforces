#include <bits/stdc++.h>
#include<stdio.h>

#include<algorithm>

#include<iostream>

#include<queue>

#define N 200005

#define LL long long

using namespace std;

inline int read( )

{

  int sum=0;char c=getchar( );bool f=0;

  while(c<'0' || c>'9') {if(c=='-') f=1;c=getchar( );}

  while(c>='0' && c<='9') {sum=sum*10+c-'0';c=getchar( );}

  if(f) return -sum;

	return sum;

}

struct ex{LL v,pos;}G;

priority_queue<ex>q;

inline bool operator < (const ex &a,const ex &b) {return a.v>b.v;}

int n,m,x;

LL a[N];

int num;

int main( )

{

	int i,j,p;

	n=read( );m=read( );x=read( );

	for(i=1;i<=n;i++)

		{

			a[i]=read( );

			if(a[i]<0) {num++;q.push((ex){-a[i],i});}

			else q.push((ex){a[i],i});

		}

	while(m--)

		{

			p=q.top( ).pos;q.pop( );

			num-=(a[p]<0);

			if(num&1) a[p]+=x;else a[p]-=x;

			num+=(a[p]<0);

			if(a[p]<0) q.push((ex){-a[p],p});

			else q.push((ex){a[p],p});

		}

	for(i=1;i<=n;i++) printf("%lld ",a[i]);

	return 0;

}