from __future__ import print_function

import json
import os
import pickle
import rnn
import string
import torch

os.environ["CUDA_VISIBLE_DEVICES"] = ""
device = torch.device("cpu")
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


def handler(_, context):
    output_names = list(rnn_model.samples(language, start_letters))

    json_region = os.environ['AWS_REGION']
    return {
        "statusCode": 200,
        "headers": {
            "Content-Type": "application/json"
        },
        "body": json.dumps({
            "Region ": json_region,
            "RequestID": context.aws_request_id,
            "TimestampChain": '[0]',
        })
    }
