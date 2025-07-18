#include <iostream>

#include <map>

#include <set>

#include <algorithm>



using namespace std;



map<string,int> mp;

map<int,string> mv;



int G[18];

int main(){



    int n,m; cin >> n >> m;

    for (int i = 0; i < n; ++i){

        string x; cin >> x;

        mp[x] = i;

        mv[i] = x;

    }

    for (int i = 0; i < m; ++i){

        string x,y;

        cin >> x >> y;

        G[mp[x]] |= (1<<mp[y]);

        G[mp[y]] |= (1<<mp[x]);

    }

    int ans = 0;

    int id;

    for (int i = 0; i < (1<<n); ++i){

        int mask = i;

        for (int j = 0; j < n; ++j){

            if(i&(1<<j))

                mask &= ((1<<n)-1)-G[j];

        }

        int c = 0;

        for (int j = 0; j < n; ++j)

            if(mask&(1<<j))c++;

        if(c > ans){

            ans = c;

            id = mask;

        }



    }

    set<string> s;

    for (int i = 0; i < n;++i)

        if(id&(1<<i)) s.insert(mv[i]);



    if(!s.empty()){

        cout << s.size() << endl;

        for(string x: s) cout << x << endl; 

    }else{

        cout << 1 << endl;

        cout << mv[0] << endl;

    }

    return 0;

}