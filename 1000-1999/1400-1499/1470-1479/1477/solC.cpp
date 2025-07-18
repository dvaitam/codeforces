#include <bits/stdc++.h>
using namespace std;
using ll = long long;

vector<ll> x, y;

bool sharp(int a, int b, int c) {
    ll dot = (x[a] - x[b]) * (x[c] - x[b]) + (y[a] - y[b]) * (y[c] - y[b]);
    return dot > 0;
}

void distribute(const vector<int>& axis, int take, vector<int>& quarter1, vector<int>& quarter2) {
    for (int i = 0; i < axis.size(); ++i) {
        if (i < take) {
            quarter1.push_back(axis[i]);
        } else {
            quarter2.push_back(axis[i]);
        }
    }
}

void append(vector<int>& line, const vector<int>& quarter1, const vector<int>& quarter2) {
    // assert(quarter1.size() == quarter2.size());
    for (int i = 0; i < quarter1.size(); ++i) {
        line.push_back(quarter1[i]);
        line.push_back(quarter2[i]);
    }
}

void append_rev(vector<int>& line, const vector<int>& quarter1, const vector<int>& quarter2) {
    // assert(quarter1.size() == quarter2.size());
    for (int i = quarter1.size() - 1; i >= 0; --i) {
        line.push_back(quarter1[i]);
        line.push_back(quarter2[i]);
    }
}

int read() {
    int n;
    cin >> n;
    x.resize(n + 1);
    y.resize(n + 1);
    for (int i = 1; i <= n; ++i) {
        cin >> x[i] >> y[i];
    }
    return n;
}

int randomize() {
    srand(time(nullptr));
    int n = 100;
    x.resize(n + 1);
    y.resize(n + 1);
    for (int i = 1; i <= n; ++i) {
        x[i] = rand() % 1000;
        y[i] = rand() % 1000;
    }
    return n;
}

int main() {
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    int n = read();
    // int n = randomize();
    bool odd = n % 2; // we want even elements, we will add the last separately if n is odd.
    if (odd) {
        --n;
    }
    vector<int> x_order(n), y_order(n);
    iota(x_order.begin(), x_order.end(), 1);
    iota(y_order.begin(), y_order.end(), 1);
    
    auto median_it = x_order.begin() + x_order.size() / 2;
    std::nth_element(x_order.begin(), median_it, x_order.end(), [](int i, int j){
        return x[i] < x[j];
    });
    ll x_median = x[*median_it];

    median_it = y_order.begin() + y_order.size() / 2;
    std::nth_element(y_order.begin(), median_it, y_order.end(), [](int i, int j){
        return y[i] < y[j];
    });
    ll y_median = y[*median_it];
    
    vector<int> place(n + 1);
    vector<vector<int>> quarter(4), axis(4);
    vector<int> center;
    for (int i = 0; i < 4; ++i) {
        quarter[i].reserve(n);
        axis[i].reserve(n);
    }
    for (int i = 1; i <= n; ++i) {
        place[i] = 3 * (1 + (y[i] < y_median) - (y[i] > y_median)) +
                        1 + (x[i] < x_median) - (x[i] > x_median);
        switch (place[i]) {
            case 0: quarter[0].push_back(i); break;
            case 1: axis[1].push_back(i); break;
            case 2: quarter[1].push_back(i); break;
            case 3: axis[0].push_back(i); break;
            case 4: center.push_back(i); break;
            case 5: axis[2].push_back(i); break;
            case 6: quarter[3].push_back(i); break;
            case 7: axis[3].push_back(i); break;
            case 8: quarter[2].push_back(i); break;
            default: break;
        }
    }

    int up_count    = quarter[0].size() + quarter[1].size() + axis[1].size();
    int down_count  = quarter[2].size() + quarter[3].size() + axis[3].size();
    int right_count = quarter[0].size() + quarter[3].size() + axis[0].size();
    int left_count  = quarter[1].size() + quarter[2].size() + axis[2].size();
    while (center.size() > 0) {
        if (axis[0].size() + axis[2].size() < abs(up_count - down_count)) {
            if (left_count > right_count) {
                axis[0].push_back(center.back());
                ++right_count;
            } else {
                axis[2].push_back(center.back());
                ++left_count;
            }
        } else {
            if (down_count > up_count) {
                axis[1].push_back(center.back());
                ++up_count;
            } else {
                axis[3].push_back(center.back());
                ++down_count;
            }
        }
        center.pop_back();
    }
    int balance = down_count - up_count + axis[0].size() - axis[2].size();
    array<int, 4> take;
    take[0] = max(0, balance / 2);
    take[2] = take[0] - balance / 2;
    int diff = quarter[0].size() + take[0] - quarter[2].size() - take[2] + axis[1].size() - axis[3].size();
    take[1] = max(0, diff);
    take[3] = take[1] - diff;

    distribute(axis[0], take[0], quarter[0], quarter[3]);
    distribute(axis[1], take[1], quarter[1], quarter[0]);
    distribute(axis[2], take[2], quarter[2], quarter[1]);
    distribute(axis[3], take[3], quarter[3], quarter[2]);

    vector<int> line;
    line.reserve(n + odd);
    if (quarter[0].size() == 0) {
        append(line, quarter[1], quarter[3]);
    } else if (quarter[1].size() == 0) {
        append(line, quarter[0], quarter[2]);
    } else {
        array<int, 4> order = {0, 2, 1, 3};
        if (!sharp(quarter[order[0]].back(), quarter[order[1]].back(), quarter[order[2]].back()) ||
            !sharp(quarter[order[1]].back(), quarter[order[2]].back(), quarter[order[3]].back())) {
            swap(order[0], order[1]);
        }
        if (!sharp(quarter[order[0]].back(), quarter[order[1]].back(), quarter[order[2]].back()) ||
            !sharp(quarter[order[1]].back(), quarter[order[2]].back(), quarter[order[3]].back())) {
            swap(order[2], order[3]);
        }
        if (!sharp(quarter[order[0]].back(), quarter[order[1]].back(), quarter[order[2]].back()) ||
            !sharp(quarter[order[1]].back(), quarter[order[2]].back(), quarter[order[3]].back())) {
            swap(order[0], order[1]);
        }
        // assert(sharp(quarter[order[0]].back(), quarter[order[1]].back(), quarter[order[2]].back()) &&
        //    sharp(quarter[order[1]].back(), quarter[order[2]].back(), quarter[order[3]].back()));
        append(line, quarter[order[0]], quarter[order[1]]);
        append_rev(line, quarter[order[2]], quarter[order[3]]);
    }
    if (odd) {
        ++n;
        line.push_back(n);
        int i = line.size() - 1;
        while (i > 1 && !sharp(line[i - 2], line[i - 1], n)) {
            line[i] = line[i - 1];
            --i;
        }
        if (i == 1 && !sharp(line[0], n, line[1])) {
            line[1] = line[0];
            --i;
        }
        line[i] = n;
    }
    for (int i : line) cout << i << " ";
    cout << endl;
    return 0;
}