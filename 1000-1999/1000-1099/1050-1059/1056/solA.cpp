#include <bits/stdc++.h>
using namespace std;

#define FasterIO ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0)

bool woke[101], tram[101];

int main(){
    FasterIO;

    int t;
    cin>>t;

    memset(tram, true, sizeof tram);

    while(t--){
        int n;
        cin>>n;

        memset(woke, false, sizeof woke);

        while(n--){
            int val;
            cin>>val;
            woke[val] = true;
        }

        for(int i=1; i<=100; i++)
            tram[i] &= woke[i];
    }

    for(int i=1; i<=100; i++)
        if(tram[i])
            cout<<i<<" ";
    cout<<endl;

    return 0;
}