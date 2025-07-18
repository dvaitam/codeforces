#include <bits/stdc++.h>
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>

#define sf scanf
#define pf printf
#define pb push_back
#define mp make_pair
#define PI ( acos(-1.0) )
#define mod 1000000007
#define maxn 100005
#define IN freopen("C.in","r",stdin)
#define OUT freopen("output.txt","w",stdout)
#define FOR(i,a,b) for(i=a ; i<=b ; i++)
#define DBG pf("Hi\n")
#define INF 1000000000000000000LL
#define i64 long long int
#define eps (1e-8)
#define xx first
#define yy second
#define ln 17
#define off 2
#define sq(x) ((x)*(x))
#define pi pair<long long,long long>
#define ds(a) ( sq(a.x)+sq(a.y) )

using namespace __gnu_pbds;
using namespace std ;

typedef tree< i64, null_type, less<i64>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;

i64 a[105][105] , r[105] , c[105] ;
i64 b[105][105] ;

int main()
{
    i64 i , j , k , l , m , n ;

    sf("%lld %lld",&n,&m) ;

    for(i=1 ; i<=n ; i++)
    {
        for(j=1 ; j<=m ; j++) sf("%lld",&a[i][j]) ;
    }

    i64 rmn = 1000 ;

    for(i=1 ; i<=m ; i++)
    {
        rmn = min( a[1][i],rmn ) ;
    }
    for(i=1 ; i<=m ; i++) c[i] = a[1][i] - rmn ;

  //  for(i=1 ; i<=m ; i++) printf("%lld ",c[i]) ;
  //  printf("\n") ;

    i64 cmn = 1000 ;

    for(i=1 ; i<=n ; i++)
    {
        cmn = min( a[i][1],cmn ) ;
    }
    for(i=1 ; i<=n ; i++) r[i] = a[i][1] - cmn ;

  //  for(i=1 ; i<=n ; i++) printf("%lld ",r[i]) ;
  //  printf("\n") ;

    i64 df = a[1][1] - r[1] - c[1] ;

    for(i=1 ; i<=n ; i++)
    {
        for(j=1 ; j<=m ; j++){
            if( (a[i][j] - r[i]-c[j])!=df ){
                printf("-1\n") ;
                return 0 ;
            }
        }
    }

    if(n<m) for(i=1 ; i<=n ; i++) r[i] += df ;
    else{
            for(i=1 ; i<=m ; i++) c[i] += df ;
    }
    i64 ans = 0 ;
    for(i=1 ; i<=n ; i++) ans += r[i] ;
    for(i=1 ; i<=m ; i++) ans += c[i] ;

    printf("%lld\n",ans) ;
    for(i=1 ; i<=n ; i++)
    {
        for(j=1 ; j<=r[i] ; j++) printf("row %lld\n",i) ;
    }
    for(i=1 ; i<=m ; i++)
    {
        for(j=1 ; j<=c[i] ; j++) printf("col %lld\n",i) ;
    }

    return 0 ;
}