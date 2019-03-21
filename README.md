[![godoc](https://godoc.org/github.com/kmiku7/go-lz4-wrapper/lz4wrapper?status.png)](https://godoc.org/github.com/kmiku7/go-lz4-wrapper/lz4wrapper)

# go-lz4-wrapper
Wrapper of go lz4 library, used to interoperable with python-lz4.

Old Library [steeve/python-lz4](https://github.com/steeve/python-lz4) compress data using lz4 block format, and append 4 bytes unsigned integer at head, which is the length of origianl data.  
There are no options to compress data in frame format, or generate a pure block.

In the new version [python-lz4](https://github.com/python-lz4/python-lz4), we can use api in lz4.block interoperate with old version. With default options, it also will append a unsigned integer at head. You can get argument descriptions at [here](https://python-lz4.readthedocs.io/en/stable/lz4.block.html#lz4.block.compress).

But there are no such function in go-binding libraries.  
In library [bkaradzic/go-lz4](https://github.com/bkaradzic/go-lz4), there seems no api to generate data in lz4 block format.  
The library [pierrec/lz4](https://github.com/pierrec/lz4) support compress data in lz4 block format, but without api to compress data with original data length appended at output string head. Even when input is not compressible, it will generate nothing.

In this repositoy, I will wrap [pierrec/lz4](https://github.com/pierrec/lz4)'s function to decompress data with a original data length at head. And also provide a compress function to deal with not-compressible data, to generate a expanded string in lz4 block format.
