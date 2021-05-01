# Copyright 2015 gRPC authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""The Python implementation of the GRPC helloworld.Greeter server."""

import grpc
import logging
import pickle
import string
import torch
from concurrent import futures

import chainfunction_pb2
import chainfunction_pb2_grpc
import rnn

torch.set_num_threads(1)

responses = ["record_response", "replay_response"]

language = 'Scottish'
language2 = 'Russian'
start_letters = 'ABCDEFGHIJKLMNOP'
start_letters2 = 'QRSTUVWXYZABCDEF'

with open('/rnn_params.pkl', 'rb') as pkl:
    params = pickle.load(pkl)

all_categories = ['French', 'Czech', 'Dutch', 'Polish', 'Scottish', 'Chinese', 'English', 'Italian', 'Portuguese',
                  'Japanese', 'German', 'Russian', 'Korean', 'Arabic', 'Greek', 'Vietnamese', 'Spanish', 'Irish']
n_categories = len(all_categories)
all_letters = string.ascii_letters + " .,;'-"
n_letters = len(all_letters) + 1

rnn_model = rnn.RNN(n_letters, 128, n_letters, all_categories, n_categories, all_letters, n_letters)
rnn_model.load_state_dict(torch.load('/rnn_model.pth'))
rnn_model.eval()


class Greeter(chainfunction_pb2_grpc.ProducerConsumerServicer):

    def InvokeNext(self, _, context):
        output_names = list(rnn_model.samples(language, start_letters))

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
