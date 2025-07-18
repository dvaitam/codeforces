#include<bits/stdc++.h>

 

#define gc getchar

#define ii(x) scanf(" %d", &x)

#define ill(x) scanf(" %lld", &x)

#define ll long long

#define pb push_back

#define mp make_pair

#define F first

#define S second

#define all(x) x.begin(),x.end()

#define fill(a,b) memset(a,b,sizeof(a))

#define rep(i,a,b) for(i=a;i<b;i++)

#define pii pair<int, int>

 

using namespace std;

 

void in(int &x){

    register int c=gc();

    x=0;

    for(;(c<48||c>57);c=gc());

    for(;c>47&&c<58;c=gc()){x=(x<<1)+(x<<3)+c-48;}

}



const int  mV = 1<<22;

int a[1000006], dp[mV + 5];



int main()

{

	int n, i, mask;

	fill(dp, -1);

	in(n);

	rep(i, 0, n){

		in(a[i]);

		dp[a[i]] = i;

	}

	for(i=1;i<=mV;i<<=1){

		rep(mask, 0, mV+1){

			if(mask & i){

				dp[mask] = max(dp[mask], dp[mask ^ i]);

			}

		}

	}

	rep(i, 0, n){

		if(dp[a[i] ^ (mV - 1)] == -1) printf("-1 ");

		else printf("%d ", a[dp[a[i] ^ (mV - 1)]]);

	}

	printf("\n");



	return 0;

}