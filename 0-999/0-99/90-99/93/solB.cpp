#include <algorithm>

#include <iostream>

#include <cstring>

#include <vector>

#include <cstdio>

#include <string>

#include <cmath>

#include <queue>

#include <set>

#include <map>

using namespace std;

typedef long long ll;

typedef double db;

typedef pair<int,int> pii;

typedef vector<int> vi;

#define de(x) cout << #x << "=" << x << endl

#define rep(i,a,b) for(int i=a;i<(b);++i)

#define all(x) (x).begin(),(x).end()

#define sz(x) (int)(x).size()

#define mp make_pair

#define pb push_back

#define fi first

#define se second

#define setIO(x) freopen(x".in","r",stdin);freopen(x".out","w",stdout)

int n , w , m;



int main(){

    scanf("%d%d%d",&n,&w,&m);

    if(n < m && m%(m-n)>0) return puts("NO") , 0;

    puts("YES");

    int cur = 1 , used = 0;

    rep(i,1,m+1){

        int sum = 0;

        for(int t=0;sum<n;++t){

            int cnt=min(n-sum,m-used);

            used+=cnt;sum+=cnt;

            if(t) putchar(' ');

            printf("%d %.16f",cur,(db)cnt/m*w);

            if(used == m) cur++,used=0;

        }

        puts("");

    }

    return 0;

}