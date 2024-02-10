import grpc
import logging
from concurrent import futures

import chainfunction_pb2
import chainfunction_pb2_grpc


class Greeter(chainfunction_pb2_grpc.ProducerConsumerServicer):

    def InvokeNext(self, _, context):
        return chainfunction_pb2.InvokeChainReply(timestampChain='[0]')


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
    chainfunction_pb2_grpc.add_ProducerConsumerServicer_to_server(Greeter(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print('serving on port 50051...')
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig()
    serve()
