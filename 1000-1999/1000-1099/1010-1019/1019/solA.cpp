#include <bits/stdc++.h>

#ifdef moskupols
    #define debug(...) fprintf(stderr, __VA_ARGS__)
    #define cdebug(...) cerr << __VA_ARGS__
#else
    #define debug(...) do {} while (false)
    #define cdebug(...) do {} while (false)
#endif

#define timestamp(x) debug("["#x"]: %.3f\n", (double)clock() / CLOCKS_PER_SEC)

#define hot(x) (x)
#define sweet(value) (value)
#define faceless

#define WHOLE(v) (v).begin(),(v).end()
#define RWHOLE(v) (v).rbegin(),(v).rend()
#define UNIQUE(v) (v).erase(unique(WHOLE(v)),(v).end())

typedef long long i64;
typedef unsigned long long ui64;
typedef long double TReal;

using namespace std;

vector<i64> cp[3005];
vector<i64> ps[3005];

struct Later {
    bool operator() (pair<int, int> a, pair<int, int> b) const {
        return cp[a.first][a.second] > cp[b.first][b.second];
    }
};

class TSolver {
public:
    int n, m;

    explicit TSolver(std::istream& in) {
        in >> n >> m;
        for (int i = 0; i < n; ++i) {
            int p, c;
            in >> p >> c;
            --p;
            cp[p].push_back(c);
        }
    }

    i64 ans = 1LL << 60;

    i64 cost(int bound) {
        i64 ret = 0;
        int tota = int(cp[0].size());

        priority_queue<pair<int, int>, vector<pair<int, int>>, Later> pq;

        for (int i = 1; i < m; ++i) {
            int here = max(0, int(cp[i].size()) - bound);
            if (here) {
                ret += ps[i][here - 1];
                tota += here;
            }
            if (here < int(cp[i].size())) {
                pq.emplace(i, here);
            }
        }

        while (tota <= bound && !pq.empty()) {
            auto top = pq.top();
            pq.pop();
            ret += cp[top.first][top.second];
            if (top.second + 1 < int(cp[top.first].size())) {
                pq.emplace(top.first, top.second + 1);
            }
            ++tota;
        }

        cdebug(bound << ' ' << tota << ' ' << ret << '\n');
        if (tota > bound) {
            return ret;
        } else {
            return 1LL << 60;
        }
    }

    void Solve() {
        for (int i = 0; i < m; ++i) {
            sort(WHOLE(cp[i]));
            partial_sum(WHOLE(cp[i]), back_inserter(ps[i]));
            assert(ps[i].size() == cp[i].size());
        }

        for (int i = 0; i <= n; ++i) {
            auto here = cost(i);
            if (here > ans) {
                break;
            }
            ans = here;
        }
    }

    void PrintAnswer(std::ostream& out) const {
        out << ans << '\n';
    }
};

int main() {
    std::ios_base::sync_with_stdio(false);
    std::cin.tie(nullptr);
    std::cout.tie(nullptr);

    {
        auto solver = std::make_shared<TSolver>(std::cin);
        solver->Solve();
        solver->PrintAnswer(std::cout);
    }

    timestamp(end);
    return 0;
}