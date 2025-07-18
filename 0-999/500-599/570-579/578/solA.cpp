#include <bits/stdc++.h>
using namespace std;

const double off = -1e-9;

int main(int argc, char** argv) {
    ios_base::sync_with_stdio(false);
#ifdef IN_FILE
    ifstream in_file(argv[1]);
    auto sbuf = cin.rdbuf(in_file.rdbuf());
#endif

    double x, y;
    cin >> x >> y;
    if (x < y) {
        cout << "-1\n";
    } else {
        cout << setprecision(12) << fixed << (x + y) / (2.  * (int)((x + y) / (2 * y))) << endl;
    }

#ifdef IN_FILE
    cin.rdbuf(sbuf);
    in_file.close();
#endif

    return 0;
}