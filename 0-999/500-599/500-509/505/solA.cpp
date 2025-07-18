#include<iostream>

#include <algorithm>

#include<string>

using namespace std;



int main(){

	string s;

	cin >> s;

	

	for (long i = 'a'; i <= 'z'; i++){

		for (long j = 0; j <= s.length(); j++){

			string t = s.substr(0, j) + string(1, i) + s.substr(j);

			string xx = t;

			reverse(t.begin(), t.end());

			if(t == xx){

				cout << t;

				return 0;

			}

		}

	}

	cout << "NA";

	

	return 0;

}