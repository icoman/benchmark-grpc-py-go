
import argparse
import logging
import time
import sys
import random
import binascii
import struct
import grpc


import benchmark_pb2
import benchmark_pb2_grpc

import mycrc32


def MkBuffer(sz):
    bytes_list = [random.randint(0,0xff) for _ in range(sz)]
    return struct.pack('B'*sz, *bytes_list)


def main():

    parser = argparse.ArgumentParser(description='Py gRPC client')
    parser.add_argument('-secured', help='Use SSL', action='store_true')
    parser.add_argument('-host', default='localhost:8080', help='address of gRPC server')
    parser.add_argument('-bs', default=300, type=int, help='Batch size in KBytes (max 4 megs)')
    parser.add_argument('-n', default=3, type=int, help='Number of tests')
    args = parser.parse_args()



    if args.secured:
        logging.info ('py RPC Client starts secured (SSL) for {}'.format(args.host))
        with open('../go-server/server.crt', 'rb') as f:
            trusted_certs = f.read()
        credentials = grpc.ssl_channel_credentials(root_certificates=trusted_certs)
        channel = grpc.secure_channel(args.host, credentials,
        options=(
            # options are defined in file
            # C:\Python27\Lib\site-packages\grpc\_cython\cygrpc.pyd
            ('grpc.ssl_target_name_override', 'MYPC'), # MYPC is hostname
        )
    )
    else:
        logging.info ('py RPC Client starts insecured (no SSL) for {}'.format(args.host))
        channel = grpc.insecure_channel(args.host)

    stub = benchmark_pb2_grpc.MyRPC1Stub(channel)
    for i in range(args.n):
        t1 = time.time()
        name="Test py client {}".format(i)
        buffer = MkBuffer(args.bs * 1024)
        bench = benchmark_pb2.Benchmark(name=name, buffer=buffer)
        crc = mycrc32.crc32(buffer)
        #crc = binascii.crc32(buffer, 0) & 0xFFFFFFFF
        response = stub.Test1(bench, timeout=10)
        t2 = time.time()
        logging.info("Response: {} in {:.2f} sec, valid: {}".format(response.text, (t2-t1), crc==response.crc32))
    channel.close()
    logging.info ("py RPC Client stops")

if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG, format='%(levelname)s %(asctime)s %(message)s', datefmt='%d-%b-%Y %H:%M:%S')
    main()
