//"Perhaps..,we 're asking the wrong questions .." -Agent Brown

#include <iostream>

#include <bits/stdc++.h>

#include <string>

#include <iomanip>

#include <unordered_map>

#define sz(s) (int(s.size()))

#define MemS(s, x) memset(s, x, sizeof(s))

#define PI acos(-1)

typedef long long ll;

using namespace std;

const ll Mod = 1e9 + 7;

const ll N = 1e4 + 1, O_O = LONG_LONG_MAX, T_T = INT_MAX, V_V = INT_MIN;

/*

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/tree_policy.hpp>

using namespace __gnu_pbds;

typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> Set;

typedef tree<int, null_type, greater_equal<int>, rb_tree_tag, tree_order_statistics_node_update> Or_Set;

*/

// 8 neighbors



int Dx[] = {-1, -1, -1, 0, 0, 1, 1, 1}; // Knight_mov:2, 1, -1, -2, -2, -1, 1, 2

int Dy[] = {-1, 0, 1, -1, 1, -1, 0, 1}; // 1, 2, 2, 1, -1, -2, -2, -1



// 4 neighbors



int dx[] = {0, 0, 1, -1};

int dy[] = {-1, 1, 0, 0};

string dir = "LRDU";

//"Look deep into your soul, into the dark and foggy mist of your memories"



void SADIEM()

{

//    freopen("input.txt", "r", stdin);

    //   freopen("output.txt", "w", stdout);

    std::ios_base::sync_with_stdio(NULL);

    cin.tie(0);

    cout.tie(0);

}



/*



 “ I’d like to let you in on a very important secret I learned when I was about your age, boy.

 You see, power, real power doesn’t come to those who were born strongest or fastest or smartest.

 No. It comes to those who will do anything to achieve it.”



*/



int main()

{

    SADIEM();

    int t;

    cin >> t;

    while (t--)

    {

        int n;

        cin >> n;

        vector<int> v(2 * n);

        for (int i = 0; i < 2 * n && cin >> v[i]; i++)

            ;

        sort(v.begin(), v.end());

        int l = 0, r = (n << 1) - 1;

        while (l < r)

            cout << v[l++] << " " << v[r--]<<" ";

        cout << (t ? "\n" : "");

    }

}