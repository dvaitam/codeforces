/*

ID: CaoLei

PROG: d.cpp

LANG: C++

*/



#include <cstdio>

#include <iostream>

#include <cstring>

#include <algorithm>

#include <set>

#include <queue>

#include <map>

#include <cmath>

#include <vector>

using namespace std;

#define N 500010

#define pi acos(-1.0)

#define inf 0x3f3f3f3f

#define pb(x) push_back((x))

typedef long long ll;

typedef unsigned long long ull;

ll dp[70][70];

ll x[70];

ll m;

int k,e;

ll dfs(int pos,int cnt,bool flag){

    if(pos==0) return cnt==k;

    if(!flag&&dp[pos][cnt]!=-1) return dp[pos][cnt];

    int u=flag?x[pos]:1;

    ll ret=0;

    for(int i=0;i<=u;i++){

        ret+=dfs(pos-1,cnt+i,i==u&&flag);

    }

    if(!flag) dp[pos][cnt]=ret;

    return ret;

}



ll ju(ll mid){

  

    e=0;

    while(mid){

        x[++e]=mid%2;

        mid>>=1;

    }

    ll tmp1=dfs(e,0,1);

    e++;

    for(int i=e;i>=2;i--){

        x[i]=x[i-1];

    }

    x[1]=0;

    ll tmp2=dfs(e,0,1);

  

    return tmp2-tmp1;

}

int main(){

    while(cin>>m>>k){

        memset(dp,-1,sizeof(dp));

        ll l=1,r=1e18+1,mid;

        while(l<=r){

         

            mid=(l+r)/2;

            ll tmp=ju(mid);

            if(tmp>m) r=mid;

            else if(tmp==m){ 

                cout<<mid<<endl;

                break;

            }

            else l=mid+1;

        }



    }

    return 0;

}