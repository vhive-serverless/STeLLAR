FROM docker.io/vhiveease/aws-python:latest
RUN pip install torch rnn

COPY server.py ./
COPY rnn_model.pth ./
COPY rnn_params.pkl ./
CMD ["server.handler"]
