#include<bits/stdc++.h>

using namespace std;



#define scl(x) scanf("%lld",&x)

#define sc(x)  scanf("%d",&x)

#define ll long long

#define lop(i,n) for(int i=0;i<n;++i)

typedef pair<int,int> ii;

typedef pair<ll,ll> pll;



int n,mat[55][55],arr[55];



int main(){

#ifndef ONLINE_JUDGE

	freopen("i.txt","r",stdin);

#endif

	sc(n);

	lop(i,n)lop(j,n)sc(mat[i+1][j+1]);

	for(int v=1;v<=n;v++){

		for(int idx=1;idx<=n;idx++){

			if(arr[idx])continue;

			int c=0;

			for(int j=1;j<=n;j++)if(mat[idx][j]==v)++c;

			if(c==n-v){

				arr[idx]=v;

				break;

			}

		}

	}

	for(int i=1;i<=n;i++)

		printf("%d ",arr[i]);

}