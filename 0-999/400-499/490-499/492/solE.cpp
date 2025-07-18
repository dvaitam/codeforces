#include <bits/stdc++.h>
#include<stdio.h>

#include<string.h>

#include<iostream>

using namespace std;

int a[1000005];

void exgcd(long long a,long long b,long long &gcd,long long &x,long long &y)

{

    if(!b)

    {

        x=1; y=0;  gcd=a;

        return;

    }

    exgcd(b,a%b,gcd,y,x);

    y-=a/b*x;

}

int main()

{

    long long n,m,dx,dy,xi,yi,x,y,gcd;

    scanf("%lld%lld%lld%lld",&n,&m,&dx,&dy);

    exgcd(n,-dx,gcd,x,y);

    y=(y*gcd%n+n)%n;

    memset(a,0,sizeof(a));

    while(m--)

    {

        scanf("%lld%lld",&xi,&yi);

        a[(yi+xi*y%n*dy%n+n)%n]++;

    }

    int i,ans=0;

    for(i=1;i<n;i++)

        if(a[ans]<a[i]) ans=i;

    printf("0 %d\n",ans);

}