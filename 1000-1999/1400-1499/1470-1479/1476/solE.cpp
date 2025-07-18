#include <bits/stdc++.h>

#define pii pair<int, int>

#define pll pair<long long, long long>

#define L(n) (n << 1)

#define R(n) (n << 1 | 1)

#define print_vector(n) for(auto a0 : n) cout << a0 << ' '; cout << endl;

#define vector_sort(n) sort(n.begin(), n.end())

#define print_array(n, l, r) for(int a0 = l; a0 <= r; a0++) cout << n[a0] << ' '; cout << '\n';

#define REP(i, l, r) for (int i = l; i <= r; i++) 

#define VREP(i, l, r) for (int i = l; i >= r; i--) 

#define MIN(a, b) (a < b ? a : b)

#define MAX(a, b) (a > b ? a : b)

#define ABS(a) ((a) > 0 ? (a) : -(a))

using namespace std;



template<class T>

istream & operator >> (istream &in, pair<T, T> &p) {

  in >> p.first >> p.second;

  return in;

}



template<class T>

ostream & operator <<(ostream &out, pair<T, T> &p) {

  out << p.first << ' ' << p.second;

  return out;

}



template<class Tuple, std::size_t N>

struct TuplePrinter {

  static void print(ostream &out, const Tuple& t) {

    TuplePrinter<Tuple, N-1>::print(out, t);

    out << ' ' << get<N-1>(t);

  }

};



template<class Tuple>

struct TuplePrinter<Tuple, 1> {

  static void print(ostream &out, const Tuple& t) {

    out << get<0>(t);

  }

};



template<class... Args>

ostream & operator <<(ostream &out, const tuple<Args...> &p) {

  TuplePrinter<decltype(p), sizeof...(Args)>::print(out, p);

  return out;

}



int N, M, K;



string A[100005];

string B[100005];

int P[100005];



map<int, int> id;

vector<int> adj[600005];

int par[600005];



int getCharId(char ch) {

  return ch == '_' ? 26 : ch - 'a';

}



int hsh(string S) {

  int res = 0;

  int mult = 1;

  for (int i = 0; i < K; i++) {

    res += mult * getCharId(S[i]);

    mult *= 27;

  }

  return res;

}



int match(string S, string T) {

  for (int i = 0; i < K; i++) {

    if (S[i] != '_' && T[i] != '_' && T[i] != S[i]) return false;

  }

  return true;

}



int isValid[600005];

int b_mx = 0;



void addEdge(int u, int v) {

  if (u == v || !isValid[u] || !isValid[v]) return;

  // cout << u << " " << v << "\n";

  adj[u].emplace_back(v);

  par[v]++;

}



void setGraph(string S, string T) {

  int u = hsh(T);

  for (int i = 0; i < b_mx; i++) {

    string tmp = S;

    for (int j = 0; j < K; j++) {

      if (i & (1 << j)) tmp[j] = '_';

    }

    // cout << i << ": " << tmp << "\n";

    addEdge(u, hsh(tmp));

  } 

}



void topoSort() {

  set<int> can;

  for (int i = 1; i <= N; i++) {

    int u = hsh(A[i]);

    if (par[u] == 0) {

      can.emplace(u);

    }

    // cout << u << " - ";

    // cout << par[u] << "\n";

  }



  vector<int> ans;

  // print_vector(can);



  while (ans.size() < N) {

    // print_vector(can);

    if (can.empty()) {

      cout << "NO\n";

      exit(0);

    }



    int u = *can.begin();

    can.erase(can.begin());

    ans.emplace_back(id[u]);



    for (int v : adj[u]) {

      par[v]--;

      if (par[v] == 0) can.emplace(v);

    }

  }



  cout << "YES\n";

  print_vector(ans);

}



int32_t main() {

  ios_base::sync_with_stdio(0);

  cin >> N >> M >> K;

  b_mx = 1 << K;  



  for (int i = 1; i <= N; i++) {

    cin >> A[i];

    int u = hsh(A[i]);

    isValid[u] = 1;

    id[u] = i;

  }



  for (int i = 1; i <= M; i++) {

    cin >> B[i] >> P[i];

    if (!match(B[i], A[P[i]])) {

      cout << "NO\n";

      return 0;

    }



    setGraph(B[i], A[P[i]]);

  }

  topoSort();

}