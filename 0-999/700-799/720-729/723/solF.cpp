#include<bits/stdc++.h>

#include<ext/pb_ds/priority_queue.hpp>

#include<ext/pb_ds/assoc_container.hpp>

#include<ext/pb_ds/tree_policy.hpp>

#include<ext/rope>



using namespace std;

using __gnu_cxx::rope;



const int oo = 0x3f3f3f3f;



#define sp putchar(' ')

#define ln putchar('\n')

#define mp make_pair

#define pb push_back

#define disp(x) cout << #x << " = " << x

#ifdef __linux__

#define getchar getchar_unlocked

#define putchar putchar_unlocked

#endif



typedef long long ll;

typedef pair<int, int> pii;

typedef __gnu_pbds::priority_queue<int, greater<int>, __gnu_pbds::pairing_heap_tag> heap;

typedef __gnu_pbds::tree<int, int, less<int>, __gnu_pbds::rb_tree_tag, __gnu_pbds::tree_order_statistics_node_update> rb_tree; //find_by_order   order_of_key (0)

template<typename T>inline bool chkmin(T &a, T b){return b < a ? a = b, true : false;}

template<typename T>inline bool chkmax(T &a, T b){return b > a ? a = b, true : false;}



template<class T> inline T read(T &x)

{

    int sign = 1;

    char c = getchar();

    for(; !isdigit(c); c = getchar())

        if(c == '-')

            sign = -1;

    for(x = 0; isdigit(c); c = getchar())

        x = x * 10 + c - '0';

    return x *= sign;

}

template<class T> inline void write(T x)

{

    if(x == 0) {putchar('0'); return;}

    if(x < 0) {putchar('-'); x = -x;}

    static char s[20];

    int top = 0;

    for(; x; x /= 10)

        s[++ top] = x % 10 + '0';

    while(top)

        putchar(s[top --]);

}



string program_name = __FILE__;

string iput = program_name.substr(0, program_name.length() - 4) + ".in";

string oput = program_name.substr(0, program_name.length() - 4) + ".out";



const int maxN = 2e5 + 100;

const int maxM = 4e5 + 100;



struct Edge {

    int u, v;

}e[maxM];

int n, m, s, t, ds, dt;

int father[maxN];



vector<pii> x, y, z, ans;



int a[maxN], b[maxN];



int find(int u)

{

    return father[u] == u ? u : father[u] = find(father[u]);

}



int main()

{

    if(fopen(iput.c_str(), "r")) {

        freopen(iput.c_str(), "r", stdin);

        freopen(oput.c_str(), "w", stdout);

    }



    read(n), read(m);

    for(int i = 1; i <= m; ++i) {

        read(e[i].u), read(e[i].v);

        if(e[i].u > e[i].v)

            swap(e[i].u, e[i].v);

    }

    read(s), read(t), read(ds), read(dt);

    if(s > t)

        swap(s, t), swap(ds, dt);



    for(int i = 1; i <= n; ++i)

        father[i] = i;



    for(int i = 1; i <= m; ++i) {

        int u = e[i].u, v = e[i].v;

        if(u != s && v != s && u != t && v != t) {

            if(find(u) != find(v)) {

                father[find(u)] = find(v);

                ans.push_back(mp(u, v));

            }

        }

    }



    bool flag = false;

    for(int i = 1; i <= m; ++i) {

        int u = e[i].u, v = e[i].v;

        if(u == s && v == t)

            flag = true;

        else {

            if(u == s) a[find(v)] = v;

            if(v == s) a[find(u)] = u;

            if(u == t) b[find(v)] = v;

            if(v == t) b[find(u)] = u;

        }

    }



    // for(int i = 1; i <= n; ++i)

    //     write(a[i]), sp; ln;

    // for(int i = 1; i <= n; ++i)

    //     write(b[i]), sp; ln;



    for(int i = 1; i <= n; ++i)

        if(find(i) == i && i != s && i != t) {

            if(a[i] && b[i]) x.push_back(mp(a[i], b[i]));

            else if(a[i]) y.push_back(mp(s, a[i]));

            else if(b[i]) z.push_back(mp(b[i], t));

            else {

                puts("No");

                return 0;

            }

        }



    // printf("%d %d %d\n", x.size(), y.size(), z.size());



    ds -= y.size(); dt -= z.size();

    for(int i = 0, _ = y.size(); i < _; ++i) ans.push_back(y[i]);

    for(int i = 0, _ = z.size(); i < _; ++i) ans.push_back(z[i]);





    // if(n >= 1000)

        // printf("%d %d\n", ds, dt);

        // printf("size = %d   flag = %d\n", (int)x.size(), flag);



    if(x.size()) {

        -- ds, -- dt;

        ans.push_back(mp(x.back().first, s));

        ans.push_back(mp(x.back().second, t));

        x.pop_back();





        // if(n >= 1000)

            // printf("%d %d\n", ds, dt);



        int tmp = max(0, min(ds, (int)x.size()));

        ds -= tmp;



        for(int i = 0; i < tmp; ++i)

            ans.push_back(mp(x[i].first, s));



        dt -= x.size() - tmp;

        for(int i = tmp, _ = x.size(); i < _; ++i)

            ans.push_back(mp(x[i].second, t));





        // if(n >= 1000)

            // printf("%d %d\n", ds, dt);

    } else if(flag) {

        ans.push_back(mp(s, t));

        -- ds, -- dt;

    } else {

        puts("No");

        return 0;

    }



    if(ds < 0 || dt < 0) {

        puts("No");

        return 0;

    }



    // if(n >= 1000)

    //     printf("%d %d\n", ds, dt);



    puts("Yes");

    for(int i = 0, _ = ans.size(); i < _; ++i)

        write(ans[i].first), sp, write(ans[i].second), ln;

    return 0;

}