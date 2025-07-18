#include <bits/stdc++.h>

using namespace std;

template<typename T>

int SIZE(T (&t)){

    return t.size();

}



template<typename T, size_t N>

int SIZE(T (&t)[N]){

    return N;

}



string to_string(char t){

    return "'" + string({t}) + "'";

}



string to_string(bool t){

    return t ? "true" : "false";

}



string to_string(const string &t, int x1=0, int x2=1e9){

    string ret = "";

    for(int i = min(x1,SIZE(t)), _i = min(x2,SIZE(t)-1); i <= _i; ++i){

        ret += t[i];

    }

    return '"' + ret + '"';

}



string to_string(const char* t){

    string ret(t);

    return to_string(ret);

}



template<size_t N>

string to_string(const bitset<N> &t, int x1=0, int x2=1e9){

    string ret = "";

    for(int i = min(x1,SIZE(t)); i <= min(x2,SIZE(t)-1); ++i){

        ret += t[i] + '0';

    }

    return to_string(ret);

}



template<typename T, typename... Coords>

string to_string(const T (&t), int x1=0, int x2=1e9, Coords... C);



template<typename T, typename S>

string to_string(const pair<T, S> &t){

    return "(" + to_string(t.first) + ", " + to_string(t.second) + ")";

}



template<typename T, typename... Coords>

string to_string(const T (&t), int x1, int x2, Coords... C){

    string ret = "[";

    x1 = min(x1, SIZE(t));

    auto e = begin(t);

    advance(e,x1);

    for(int i = x1, _i = min(x2,SIZE(t)-1); i <= _i; ++i){

        ret += to_string(*e, C...) + (i != _i ? ", " : "");

        e = next(e);

    }

    return ret + "]";

}



template<int Index, typename... Ts>

struct print_tuple{

    string operator() (const tuple<Ts...>& t) {

        string ret = print_tuple<Index - 1, Ts...>{}(t);

        ret += (Index ? ", " : "");

        return ret + to_string(get<Index>(t));

    }

};



template<typename... Ts>

struct print_tuple<0, Ts...> {

    string operator() (const tuple<Ts...>& t) {

        return to_string(get<0>(t));

    }

};



template<typename... Ts>

string to_string(const tuple<Ts...>& t) {

    const auto Size = tuple_size<tuple<Ts...>>::value;

    return print_tuple<Size - 1, Ts...>{}(t);

}



void dbgr(){;}

template<typename Heads, typename... Tails>

void dbgr(Heads H, Tails... T){

    cout << to_string(H) << " | ";

    dbgr(T...);

}



void dbgs(){;}

template<typename Heads, typename... Tails>

void dbgs(Heads H, Tails... T){

    cout << H << " ";

    dbgs(T...);

}



/*

formatted functions:

*/



/*

consider __VA_ARGS__ as a whole:

dbgv() prints values only

dbg() prints name and values

*/

#define dbgv(...) cout << to_string(__VA_ARGS__) << endl;



#define dbg(...) cout << "[" << #__VA_ARGS__ << "]: "; dbgv(__VA_ARGS__);

//#define dbg(...)



/*

consider __VA_ARGS__ as a sequence of arguments:

dbgr() prints values only

dbgm() prints names and values

*/

#define dbgr(...) dbgr(__VA_ARGS__); cout << endl;



#define dbgm(...) cout << "[" << #__VA_ARGS__ << "]: "; dbgr(__VA_ARGS__);



struct custom_hash {

    static uint64_t splitmix64(uint64_t x) {

        // http://xorshift.di.unimi.it/splitmix64.c

        x += 0x9e3779b97f4a7c15;

        x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9;

        x = (x ^ (x >> 27)) * 0x94d049bb133111eb;

        return x ^ (x >> 31);

    }



    size_t operator()(uint64_t x) const {

        static const uint64_t FIXED_RANDOM = chrono::steady_clock::now().time_since_epoch().count();

        return splitmix64(x + FIXED_RANDOM);

    }

};



#define int long long

typedef pair<int, int> pi;

const int mod=1000000007;

const int maxN=200001;

const double eps=1e-11;

string s;

int n, m, k;

void solve() {

    cin>>n>>m;

    vector<pi> dishes(n), out(n);

    int eat=n*m, fish=0, meat=0;

    int minfish=0, maxfish=0, minmeat=0, maxmeat=0, best=LLONG_MAX;

    for (int i=0; i<n; i++) {

        int a, b;

        cin>>a>>b;

        fish+=a, meat+=b;

        int reqfish=max(m-b, 0LL), reqmeat=max(m-a, 0LL);

        maxfish+=min(m, a), maxmeat+=min(m, b), minfish+=reqfish, minmeat+=reqmeat;

        out[i].first=reqfish, out[i].second=reqmeat;

        dishes[i]={a, b};

    }



    fish-=minfish, meat-=minmeat;

    int both=maxfish-minfish, test=maxmeat-minmeat;

    assert(both==test);

    //dbgm(fish, meat, minfish, maxfish, minmeat, maxmeat, both);

    int diff=abs(fish-meat), eatfish=0, eatmeat=0;

    if (diff>both) {

        best=diff-both;

        if (fish>meat) eatfish=both;

        else eatmeat=both;

    }



    else {

        best=(both-diff)%2;

        if (fish>meat) eatfish=diff;

        else eatmeat=diff;

        eatfish+=(both-diff)/2, eatmeat+=(both-diff)-(both-diff)/2;

    }



    //dbgm(eatfish, eatmeat);

    for (int i=0; i<n; i++) {

        int eat=m-out[i].first-out[i].second;

        int mx=min(eatfish, eat);

        eat-=mx, eatfish-=mx, out[i].first+=mx;

        eatmeat+=eat, out[i].second+=eat;

    }



    cout<<best<<"\n";

    for (auto p:out) cout<<p.first<<" "<<p.second<<"\n";

    //IMPORTANT!!!! Remember to change back the value of maxN after debugging!!!!!

    //Take all of the inputs before returning!!!!!

    //Open the topics list!!!!!

}



int32_t main() {

    std::ios::sync_with_stdio(false);

    cin.tie(NULL);

    int t;

    cin>>t;

    while (t--) {

        solve();

    }

}