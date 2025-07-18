#include <cstdio>
#include <cstring>
#include <algorithm>
using namespace std;
typedef long long ll;
int main(){
    ll n,k;scanf("%lld%lld",&n,&k);
    if(n>k*(k-1)) printf("NO\n");
    else{
        printf("YES\n");
        ll cnt=0,add=0,fi=0,se;
        while(cnt<n){
            if(cnt%k==0) add++;
            cnt++;
            fi=fi%k+1ll;
            se=(fi+add-1ll)%k+1ll;
            printf("%lld %lld\n",fi,se);
        }
    }
    return 0;
}