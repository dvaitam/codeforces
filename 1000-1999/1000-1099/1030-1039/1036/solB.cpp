#include <bits/stdc++.h>
using namespace std;
typedef unsigned long long ll;
typedef pair<ll,ll> P;

int main() {
    int q;
    scanf("%d" ,&q);
    while(q--) {
        ll a,b,c;
        scanf("%lld%lld%lld",&a,&b,&c);
        if (a > b)
            swap(a, b);
        ll Min = a + (b-a);
        if(Min > c) {
            printf("-1\n");
            continue;
        }
        if(a!=b) {
            ll ans = a;
            ll cha = b-a;
            if(cha&1) {
                if((c-a)%2 == 0) {
                    ans += c-a-1;
                } else {
                    ans += c-a-1;
                }
            } else {
                if((c-a)%2 == 0) {
                    ans += c-a;
                } else {
                    ans += c-a-2;
                }
            }
            printf("%lld\n",ans);
        } else {
            ll ans;
            if(a == c) {
                ans = a;
                printf("%lld\n",ans);
                continue;
            } else {
                if((c-a)%2 == 0) {
                    ans = c;
                } else {
                    ans = c-2;
                }
                printf("%lld\n",ans);
            }
        }
    }
    return 0;
}