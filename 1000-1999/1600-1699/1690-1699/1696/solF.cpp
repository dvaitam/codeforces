#include <bits/stdc++.h>

#include <iostream>

using namespace std;



bool get_graph_by_edge(pair<int, int> edge, const vector<vector<vector<bool>>> &d,

                       vector<int> &visited, vector<pair<int, int>> &edges,

                       vector<vector<int>> &paths) {

    int n = d.size();

    queue<pair<int, int>> q;

    paths[edge.first][edge.second] = 1;

    paths[edge.second][edge.first] = 1;

    q.push(edge);

    while (!q.empty()) {

        auto [l, r] = q.front();

//        cout << l << ' ' << r << '\n';

        q.pop();

        for (int i = 0; i < n; ++i) {

            if (d[l][r][i]) {

//                cout << i << "!" << '\n';

                if (visited[i] and i != l) {

                    return false;

                }

                q.push({r, i});

                visited[i] = true;

                edges.push_back({r, i});

                for (int j = 0; j < n; ++j) {

                    if (paths[j][r] != -1 and j != i) {

//                        cout << i << ' ' << j << ' ' << l << ' ' << r << ' ' << paths[j][r] << '\n';

                        paths[j][i] = paths[j][r] + 1;

                        paths[i][j] = paths[j][r] + 1;

                     }

                }

            }

        }

    }

    return true;

}



bool check_tree(const vector<vector<vector<bool>>> &d, const vector<vector<int>> &paths) {

    int n = d.size();

    for (int i = 0; i < n; ++i) {

        for (int j = 0; j < n; ++j) {

            for (int k = i + 1; k < n; ++k) {

                if (d[i][j][k] == 1 and paths[i][j] != paths[j][k]) {

                    return false;

                }

                if (d[i][j][k] == 0 and paths[i][j] == paths[j][k]) {

                    return false;

                }

            }

        }

    }

    return true;

}



int main() {

    ios_base::sync_with_stdio(false);

    cin.tie(nullptr);



    int n, N;

    cin >> N;

    for (int t = 0; t < N; ++t) {

        cin >> n;

        string s;

        vector<vector<vector<bool>>> d(n, vector<vector<bool>>(n, vector<bool>(n, false)));



        for (int i = 0; i < n-1; ++i) {

            for (int j = i; j < n-1; ++j) {

                cin >> s;

                for (int k = 0; k < s.size(); ++k) {

                    if (s[k] == '1') {

                        d[i][k][j+1] = true;

                        d[j+1][k][i] = true;

                    }

                }

            }

        }

        bool found = false;

        for (int i = 1; i < n; ++i) {

//            cout << i << '\n';

            vector<pair<int, int>> edges;

            vector<int> visited(n);

            visited[0] = true;

            visited[i] = true;

            vector<vector<int>> paths(n, vector<int>(n, -1));

            for (int j = 0; j < n; ++j) {

                paths[j][j] = 0;

            }

            if (get_graph_by_edge({0, i}, d, visited, edges, paths)) {

//                cout << "ed_1!\n";

                if (get_graph_by_edge({i, 0}, d, visited, edges, paths)) {

//                    cout << "ed_2!\n";

                    edges.push_back({0, i});

                    if (edges.size() != n-1) {

                        continue;

                    }

                    if (!check_tree(d, paths)) {

                        continue;

                    }

                    cout << "Yes" <<'\n';

                    found = true;

                    for (auto [l, r] : edges) {

                        cout << l + 1 << ' ' << r + 1 <<'\n';

                    }

//                    for (int j = 0; j < n; ++j) {

//                        for (int k = 0; k < n; ++k) {

//                            cout << paths[j][k] << ' ';

//                        }

//                        cout << '\n';

//                    }

                    break;



                }

            }

        }

        if (!found) {

            cout << "No\n";

        }

    }

    return 0;

}