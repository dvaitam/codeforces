#include <bits/stdc++.h>
#define ll long long
using namespace std;


void testcase(){
    int n , k ;
    cin >> n >> k;
    pair<int,int> a[n];
    for (int i = 0 ; i < n ; i++) {
        int yy;
        cin >> yy ;
        a[i] = {yy,i+1};
    }
    sort(a,a+n);
    int minn = 2000000000, n1 , n2, i1,i2;
    for (int i = 1 ; i < n ; i++) {
        int ttt = (a[i].first ^ a[i-1].first);
        if (ttt < minn) {
            minn = ttt;
            n1 = a[i].first;
            n2 = a[i-1].first;
            i1 = a[i].second;
            i2 = a[i-1].second;
        }
    }
    int pp = 1;
    for (int i = 0 ; i < k ; i++) {
        pp *= 2;
    }
    pp--;
    cout << i1 << " " << i2 << " " << (pp^n1) << "\n";

}


int main(){
    ios::sync_with_stdio(0);
    cin.tie(0);
    int q = 1;
    cin >> q;
    while(q--){
        testcase();
    }
    return 0;
}