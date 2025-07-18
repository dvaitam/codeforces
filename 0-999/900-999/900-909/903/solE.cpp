#include <iostream>
#include <string>
#include <unordered_set>
#include <vector>
#include <algorithm>

using namespace std;

bool possible(const string& target, const string& s, bool canStay) {
    vector<int> diffs;
    for(int i=0; i<target.size(); ++i) {
        if(target[i]!=s[i])
            diffs.push_back(i);
    }

    if(canStay && diffs.empty()) return true;
    if(diffs.size()!=2) return false;
    if(target[diffs[0]]!=s[diffs[1]] || target[diffs[1]]!=s[diffs[0]])
        return false;
    return true;
}

int main() {
    int k,n;
    cin >> k >> n;

    unordered_set<string> strs;
    for(int i=0; i<k; ++i) {
        string tmp;
        cin >> tmp;
        strs.insert(tmp);
    }

    if(strs.size()==1) {
        string res = *strs.begin();
        swap(res[0],res[1]);
        cout << res;
        return 0;
    }

    bool canStay = false;
    vector<int> chrCnt(26);
    for(char c: *strs.begin()) {
        ++chrCnt[c-'a'];
        if(chrCnt[c-'a']>1)
            canStay = true;
    }

    for(const string& s: strs) {
        vector<int> cnt(26);
        for(char c: s)
            ++cnt[c-'a'];

        if(cnt != chrCnt) {
            cout << -1;
            return 0;
        }
    }

    vector<string> strVec(strs.begin(), strs.end());

    vector<int> diffs;
    for(int i=0; i<strVec[0].size(); ++i)
        if(strVec[0][i]!=strVec[1][i])
            diffs.push_back(i);

    if(diffs.size()>4 || diffs.size()<2) {
        cout << -1;
        return 0;
    }

    unordered_set<string> targets;

    if(canStay)
        targets.insert(strVec[0]);

    for(int i=0; i<diffs.size(); ++i)
    for(int j=i+1; j<diffs.size(); ++j) {
        swap(strVec[0][diffs[i]], strVec[0][diffs[j]]);
        targets.insert(strVec[0]);
        swap(strVec[0][diffs[i]], strVec[0][diffs[j]]);
    }

    for(int i=1; i<strVec.size(); ++i) {
        vector<string> toErase;
        for(const string& s: targets)
            if(!possible(s,strVec[i],canStay))
                toErase.push_back(s);
        for(string& s: toErase)
            targets.erase(s);
    }

    if(targets.size())
        cout << *targets.begin();
    else
        cout << -1;

    return 0;

}