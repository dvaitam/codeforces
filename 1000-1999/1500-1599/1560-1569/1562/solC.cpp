// ~BhupinderJ

#include <bits/stdc++.h>

using namespace std;

typedef long long ll;

#define endl "\n"

#define spc <<" "<<



int findZero(string &s, int n){

    int idx = -1;

    for(int i=0 ; i<n ; i++){

        if(s[i] == '0') idx = i+1;

    }

    return idx;

}

int main(){

    ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);

    int T; cin >> T;

    while(T--){

        int n; cin >> n;

        string s; cin >> s;

        int idx = findZero(s, n);

        if(idx > n/2) cout << 1 spc idx spc 1 spc idx-1 << endl;

        else if(idx < n/2) cout << n/2+1 spc n spc n/2 spc n-1 << endl;

        else cout << n/2+1 spc n spc n/2 spc n << endl;

    }

}

/*



*/