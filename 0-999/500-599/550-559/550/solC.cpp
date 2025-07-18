#include <iostream>

#include <string>

#include <iomanip>

#include <algorithm>

#include <vector>

using namespace std;

vector<int>vec;

void generate(){

	for (int i = 0; i < 1000; i+=8){

		vec.push_back(i);

	}

}

int main(){

	generate();

	string s;

	cin >> s;

	for (int i = 0; i < vec.size(); i++){

		string s2 = "";

		int x = vec[i];

		vector<int>indeces;

		if (x == 0)indeces.push_back(x);



		while (x){

			indeces.push_back(x % 10);

			x /= 10;

		}

		reverse(indeces.begin(), indeces.end());

		for (int j = 0; j < indeces.size(); j++){

			s2 += indeces[j] + '0';

		}

		indeces.clear();

		int k = 0;

		string s3 = "";

		for (int j = 0; j < s.size(); j++){

			if (s[j] == s2[k])k++, s3 += s[j];

		}

		if (s3 == s2)return cout << "YES" << endl << s3 << endl, 0;

	}

	cout << "NO" << endl;

	return 0;

}