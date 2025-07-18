#include<bits/stdc++.h>

using namespace std;



#define scl(x) scanf("%lld",&x)

#define sc(x)  scanf("%d",&x)

#define ll long long

#define lop(i,n) for(int i=0;i<n;++i)

typedef pair<int,int> ii;

typedef pair<ll,ll> pll;



int n,arr[1010],st,init;



bool ok(int m){

	for(int k=0;k<=2e4;k++){

		int mn=arr[0]-m;

		int mx=arr[0]+m;

		bool f=1;

		for(int i=1;i<n;i++){

			mn+=k;

			mx+=k;

			int mxa=arr[i]+m;

			int mna=arr[i]-m;

			mx=min(mx,mxa);

			mn=max(mn,mna);

			if(mn>mx){

				f=0;

				break;

			}

		}

		if(f){

			st=k;

			init=mn-(n-1)*k;

			return 1;

		}

	}

	return 0;

}

int main(){

#ifndef ONLINE_JUDGE

	freopen("i.txt","r",stdin);

#endif

	sc(n);

	lop(i,n)sc(arr[i]);

	sort(arr,arr+n);

	int s=0,e=1e4,m,bst=-1;

	while(s<=e){

		m=s+((e-s)>>1);

		if(ok(m))e=m-1,bst=m;

		else s=m+1;

	}

	printf("%d\n%d %d",bst,init,st);

}