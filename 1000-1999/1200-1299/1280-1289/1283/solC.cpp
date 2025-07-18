#include <bits/stdc++.h>
using namespace std;

int main() {
    ios::sync_with_stdio(false);
    cin.tie(0); cout.tie(0);
//↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
    int n; cin >> n;
    vector<int> v(n);
    vector<bool> used(n);
    for (int i = 0; i < n; ++i) {
        cin >> v[i];
        if (v[i])
            used[v[i] - 1] = true; // Вектор использованых людей
    }
    
    //
    deque<int> dq;
    for (int i = 0; i < n; ++i)
        if (!used[i])
            dq.push_back(i);
    
    /*cout << "DQ: ";
    for (auto it : dq)
        cout << it << " ";
    cout << endl;*/
    
    int temp = 0;
    for (auto it : dq) {
        if (!v[it]) {
            //cout << "!" << it << endl;
            if (it == dq.front()) {
                temp = dq.back();
                dq.pop_back();
            }
            else {
                temp = dq.front();
                dq.pop_front();
            }
            
            v[it] = temp + 1;
        }
    }
    
    for (auto it : v) {
        if (!it) {
            temp = dq.front();
            dq.pop_front();
            cout << temp + 1 << " ";
        }
        else cout << it << " ";
    }
//↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑
    return 0;
}