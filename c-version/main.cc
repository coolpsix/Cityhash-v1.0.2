#include <city.h>
#include <cstdio>
#include <iostream>
#include <string.h>

static const uint64_t k0 = 0xc3a5c85c97cb3127ULL;
static const int kDataSize = 1 << 20;
static const int kTestSize = 300;

static char data[kDataSize];

void setup() {
  uint64_t a = 9;
	uint64_t b = 777;
  for (int i = 0; i < kDataSize; i++) {
    a += b;
    b += a;
    a = (a ^ (a >> 41)) * k0;
    b = (b ^ (b >> 41)) * k0 + i;
    unsigned char u = b >> 37;
    memcpy(data + i, &u, 1);  // uint8 -> char
  }
}

int main(int argc, char **argv) {
	setup();
	std::cout << "var expected = []Uint128{" << std::endl;
	for(int i = 0; i < kTestSize -1; i++) {
		int offset = i*i;
		int len = i;
		auto result = CityHash_v1_0_2::CityHash128(data + offset, i);
		std::cout << "\tUint128{First: " << result.first << ", Second: " << result.second << "},";
		std::cout << "\t// iteration=" << i << " offset=" << offset << " len=" << len << std::endl;
	}
	auto result = CityHash_v1_0_2::CityHash128(data, kDataSize);
	std::cout << "\tUint128{First: " << result.first << ", Second: " << result.second << "}, // whole iter" << std::endl;
	std::cout << "}" << std::endl;
	return 0;
}
