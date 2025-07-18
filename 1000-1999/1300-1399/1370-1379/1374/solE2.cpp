#include<iostream>

#include<cstring>

#include<cassert>

#include<vector>

#include<algorithm>

using namespace std;

using LL = long long;



int main(){



#ifdef LOCAL

    freopen("data.in", "r", stdin);

    freopen("data.out", "w", stdout);

#endif



    cin.tie(0);

    cout.tie(0);

    ios::sync_with_stdio(0);



    int n, m, k;

    cin >> n >> m >> k;



    vector<pair<int, int> > v[4];

    for(int i = 1; i <= n; i++){

        int a, b, c;

        cin >> a >> b >> c;

        int state = (b << 0) | (c << 1);

        v[state].push_back({a, i});

    }

    for(int i = 0; i < 4; i++)

        sort(v[i].begin(), v[i].end());



    auto check = [&](int c){

        if (c + (k - c) + (k - c) > m) return false;

        if (k - c > (int)v[1].size()) return false;

        if (k - c > (int)v[2].size()) return false;

        if (c + v[0].size() + v[1].size() + v[2].size() < m) return false;

        return true;

    };



    if (!check(min(m, int(v[3].size())))){

        cout << -1 << '\n';

        return 0;

    }



    int pt[4] = {0};

    int sum = 0;

    for(int i = 0; i < min(m, int(v[3].size())); i++)

        sum += v[3][pt[3]++].first;



    while(pt[1] + pt[3] < k)

        sum += v[1][pt[1]++].first;



    while(pt[2] + pt[3] < k)

        sum += v[2][pt[2]++].first;



    auto s = [&](){

        return pt[0] + pt[1] + pt[2] + pt[3];

    };



    while(s() < m){

        int val = 1e5, t = -1;

        for(int i = 0; i <= 2; i++){

            if (pt[i] == v[i].size()) continue;

            if (val > v[i][pt[i]].first) t = i, val = v[i][pt[i]].first;

        }

        sum += val;

        pt[t]++;

    }



    int ans = sum;

    int res[4];

    memcpy(res, pt, sizeof res);



    // assert(s() == m);

    // assert(pt[1] + pt[3] >= k);

    // assert(pt[1] + pt[3] >= k);



    while(pt[3]){

        sum -= v[3][--pt[3]].first;

        if (!check(pt[3])) break;

        while(pt[1] + pt[3] < k)

            sum += v[1][pt[1]++].first;



        while(pt[2] + pt[3] < k)

            sum += v[2][pt[2]++].first;

        

        while(s() > m){

            sum -= v[0][--pt[0]].first;

        }



        while(s() < m){

            int val = 1e5, t = -1;

            for(int i = 0; i <= 2; i++){

                if (pt[i] == v[i].size()) continue;

                if (val > v[i][pt[i]].first) t = i, val = v[i][pt[i]].first;

            }

            // assert(t != -1);

            sum += val;

            pt[t]++;

        }



        // assert(s() == m);

        // assert(pt[1] + pt[3] >= k);

        // assert(pt[2] + pt[3] >= k);



        if (sum < ans){

            ans = sum;

            memcpy(res, pt, sizeof res);

        }

    }



    // assert(res[0] + res[1] + res[2] + res[3] == m);

    // assert(res[1] + res[3] >= k);

    // assert(res[2] + res[3] >= k);



    cout << ans << '\n';

    for(int i = 0; i < 4; i++){

        for(int j = 0; j < res[i]; j++){

            cout << v[i][j].second << ' ';

        }

    }

    cout << '\n';



}