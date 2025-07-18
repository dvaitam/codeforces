#include <bits/stdc++.h>
using namespace std;
typedef vector<int> vi;
#define sz(x) (int)(x).size()

const int MAXN = 2010;
char s[MAXN], t[MAXN];

const string drive_prefix = "DRIVE ";
const string pickup = "PICKUP", dropoff = "DROPOFF";
string drive (int x) {
    return drive_prefix + to_string(x + 1);
}
void solve (vi nom, vi noc, vi mc, vi cm, vi mno, vi cno) {
    if (sz(mc) + sz(cm) + sz(mno) + sz(cno) == 0) {
        cout << 0 << endl;
        exit(0);
    }

    // doing nom first
    if (nom.empty()) return;
    vector<string> ans;
    // nom mc cm mc cm ...
    ans.push_back(drive(nom.back())); nom.pop_back();
    bool mcnext = true, last_m = true;
    for (;;) {
        if (mcnext) {
            if (mc.empty()) break;
            ans.push_back(pickup);
            ans.push_back(drive(mc.back()));
            mc.pop_back();
            ans.push_back(dropoff);
            last_m = false;
            mcnext = false;
        } else {
            if (cm.empty()) break;
            ans.push_back(pickup);
            ans.push_back(drive(cm.back()));
            cm.pop_back();
            ans.push_back(dropoff);
            last_m = true;
            mcnext = true;
        }
    }

    if (last_m) {
        // I can fix one mno
        if (!mno.empty()) {
            ans.push_back(pickup);
            ans.push_back(drive(mno.back()));
            mno.pop_back();
            ans.push_back(dropoff);
        }
    } else {
        // I can fix one cno
        if (!cno.empty()) {
            ans.push_back(pickup);
            ans.push_back(drive(cno.back()));
            cno.pop_back();
            ans.push_back(dropoff);
        }
    }

    // either MC = 0 or CM = 0
    assert (sz(mc) * sz(cm) == 0);

    int mcc0 = min(sz(mc), sz(cno));
    int cmm0 = min(sz(cm), sz(mno));

    for (int _ = 0; _ < cmm0; _++) {
        if (noc.empty()) return;
        ans.push_back(drive(noc.back()));
        noc.pop_back();
        ans.push_back(pickup);
        ans.push_back(drive(cm.back()));
        cm.pop_back();
        ans.push_back(dropoff);
        ans.push_back(pickup);
        ans.push_back(drive(mno.back()));
        mno.pop_back();
        ans.push_back(dropoff);
    }

    for (int _ = 0; _ < mcc0; _++) {
        if (nom.empty()) return;
        ans.push_back(drive(nom.back()));
        nom.pop_back();
        ans.push_back(pickup);
        ans.push_back(drive(mc.back()));
        mc.pop_back();
        ans.push_back(dropoff);
        ans.push_back(pickup);
        ans.push_back(drive(cno.back()));
        cno.pop_back();
        ans.push_back(dropoff);
    }

    // now, MC == M0, CM == C0
    while (!mc.empty()) { mno.push_back(mc.back()); mc.pop_back(); }
    while (!cm.empty()) { cno.push_back(cm.back()); cm.pop_back(); }

    while (!mno.empty()) {
        if (nom.empty()) return;
        ans.push_back(drive(nom.back()));
        nom.pop_back();
        ans.push_back(pickup);
        ans.push_back(drive(mno.back()));
        mno.pop_back();
        ans.push_back(dropoff);
    }

    while (!cno.empty()) {
        if (noc.empty()) return;
        ans.push_back(drive(noc.back()));
        noc.pop_back();
        ans.push_back(pickup);
        ans.push_back(drive(cno.back()));
        cno.pop_back();
        ans.push_back(dropoff);
    }

    assert (sz(mc) + sz(cm) + sz(mno) + sz(cno) == 0);

    cout << sz(ans) << endl;
    for (auto str : ans) cout << str << endl;
    exit(0);
}

int main() {
    int n; scanf("%d", &n);
    scanf("%s", s);
    scanf("%s", t);

    vector<int> nom, noc, mc, cm, mno, cno;
    for (int i = 0; i < n; i++) {
        if (s[i] == t[i]) continue;
        if (s[i] == '-') {
            if (t[i] == 'M') nom.push_back(i);
            else if (t[i] == 'C') noc.push_back(i);
            else assert(false);
        } else if (s[i] == 'C') {
            if (t[i] == 'M') cm.push_back(i);
            else if (t[i] == '-') cno.push_back(i);
            else assert(false);
        } else if (s[i] == 'M') {
            if (t[i] == 'C') mc.push_back(i);
            else if (t[i] == '-') mno.push_back(i);
            else assert(false);
        } else assert(false);
    }

    solve(nom, noc, mc, cm, mno, cno);
    solve(noc, nom, cm, mc, cno, mno);

    return 0;
}