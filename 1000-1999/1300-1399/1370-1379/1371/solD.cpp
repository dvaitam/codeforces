#include <iostream>

#include <algorithm>

#include <cstring>

#include <string>

#include <map>

#include <set>

#include <queue>

#include <stack>

#include <vector>

#include <unordered_map>

#include <unordered_set>

#include <climits>

#include <cmath>

#include <functional>

#include <numeric>

using namespace std;

using LL = long long;

using PII = pair<int, int>;

using PLL = pair<LL, LL>;

using vi = vector<int>;

using vvi = vector<vi>;

using vl = vector<LL>;

using vvl = vector<vl>;

using vb = vector<bool>;

using vvb = vector<vb>;

#define ph emplace_back

#define all(x) x.begin(), x.end()

template <typename T> bool chkMax(T &x, T y) { return (y > x) ? x = y, 1 : 0; }

template <typename T> bool chkMin(T &x, T y) { return (y < x) ? x = y, 1 : 0; }



void func()

{

    int n, k;

    cin>>n>>k;

    if(k%n)cout<<2<<'\n';

    else cout<<0<<'\n';

    vector<vector<char>>w(n, vector<char>(n, '0'));

    for(int i=0, x=0, y=0; i<k; i++){

        w[x][y] = '1';

        if(++x==n)x = 0;

        if(++y==n)y=0, x++;

    }

    for(auto& x: w){

        for(auto& y: x)cout<<y;

        cout<<'\n';

    }

}



int main(void){

    ios::sync_with_stdio(false);

    cin.tie(0);

    int t;

    cin>>t;

    while(t--)func();

    return 0;

}