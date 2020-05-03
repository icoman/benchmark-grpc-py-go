
import argparse
import logging
from concurrent import futures
import time
import grpc
import binascii
import benchmark_pb2
import benchmark_pb2_grpc

import mycrc32

class MyServer(benchmark_pb2_grpc.MyRPC1Servicer):
    """Provides methods that implement functionality of RPC server."""

    def __init__(self):
        logging.info('py RPC server init')
 

    def Test1(self, request, context):
        logging.info('Test: {} ({:.2f} KBytes)'.format(request.name, float(len(request.buffer))/1024))
        crc = mycrc32.crc32(request.buffer)
        #crc = binascii.crc32(request.buffer, 0) & 0xFFFFFFFF
        return benchmark_pb2.Response(text='py Server says text ok', crc32=crc)

def main():

    parser = argparse.ArgumentParser(description='Py gRPC client')
    parser.add_argument('-secured', help='Use SSL', action='store_true')
    parser.add_argument('-host', default='localhost:8080', help='address:port of gRPC server')
    args = parser.parse_args()


    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10), options=())
    benchmark_pb2_grpc.add_MyRPC1Servicer_to_server(MyServer(), server)
    if args.secured:
        logging.info('py gRPC server starts secured (SSL) on {}'.format(args.host))  
        # read in key and certificate
        with open('../go-server/server.key', 'rb') as f:
            private_key = f.read()
        with open('../go-server/server.crt', 'rb') as f:
            certificate_chain = f.read()

        # create server credentials
        server_credentials = grpc.ssl_server_credentials(((private_key, certificate_chain,),))

        # add secure port using crendentials
        server.add_secure_port(args.host, server_credentials)
    else:
        logging.info('py gRPC server starts insecured (no SSL) on {}'.format(args.host))  
        server.add_insecure_port(args.host)
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG, format='%(levelname)s %(asctime)s %(message)s', datefmt='%d-%b-%Y %H:%M:%S')
    main()

