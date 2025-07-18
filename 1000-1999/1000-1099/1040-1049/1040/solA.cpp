#include <bits/stdc++.h>

#define FOR(i,a,b) for(int i = a; i <= b; ++i)
#define FORD(i,a,b) for(int i = a; i >= b; --i)
#define REP(i,a) for(int i = 0; i < a; ++i)

using namespace std;

template<typename T> inline void Mini(T &x, T y){if(x > y) x = y;}
template<typename T> inline void Maxi(T &x, T y){if(x < y) x = y;}

const int N = 22;

int n, D[2], A[N];
int res;

int main(){
    ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);
    cin>>n>>D[0]>>D[1];
    FOR(i,1,n) cin>>A[i];
    FOR(i,1,n/2){
        if(A[i]!=2&&A[n-i+1]!=2){
            if(A[i]!=A[n-i+1]){
                cout<<-1<<'\n';
                return 0;
            }
            continue;
        }
        if(A[i]==2&&A[i]==A[n-i+1]){
            res += 2*min(D[0],D[1]);
            continue;
        }
        if(A[i]==2){
            res += D[A[n-i+1]];
            continue;
        }
        if(A[n-i+1]==2){
            res += D[A[i]];
            continue;
        }
    }
    if(n%2==1&&A[n/2+1]==2) res += min(D[0],D[1]);
    cout<<res<<'\n';
    return 0;
}