#include <iostream>
#include <vector>
#include <algorithm>
#include <cmath>

int main() {
    std::ios_base::sync_with_stdio(false);
    std::cin.tie(0);
    int n;
    std::cin >> n;
    std::string s;
    std::cin >> s;
    int id = 0;
    int hmin = 0;
    int h = 0;
    int count = 0;

    for (int i = 0; i < n; ++i) {
        char c = s[i];

        if (h < hmin) {
            id = i;
            hmin = h;
            count = 0;
        }

        if (h == hmin) {
            ++count;
        }

        if (c == '(') {
            ++h;
        } else {
            --h;
        }
    }

    if (h != 0) {
        std::cout << "0\n1 1\n";
        return 0;
    }

    std::string s2 = s.substr(id) + s.substr(0, id);
    int best = count;
    int curr1 = 0;
    int curr2 = 0;
    int a = 0;
    int b = 0;
    int a1 = 0, a2;
    h = 0;

    for (int i = 0; i < n; ++i) {
        char c = s2[i];

        if (c == '(') {
            ++h;
        } else {
            --h;
        }

        if (h == 0) {
            if (curr1 > best) {
                best = curr1;
                a = a1;
                b = i;
            }

            curr1 = 0;
            a1 = i + 1;
        } else if (h == 1) {
            ++curr1;

            if (curr2 + count > best) {
                best = curr2 + count;
                a = a2;
                b = i;
            }

            curr2 = 0;
            a2 = i + 1;
        } else if (h == 2) {
            ++curr2;
        }
    }

    std::cout << best << "\n";
    std::cout << (a + id) % n + 1 << " " << (b + id) % n + 1 << "\n";
    return 0;
}