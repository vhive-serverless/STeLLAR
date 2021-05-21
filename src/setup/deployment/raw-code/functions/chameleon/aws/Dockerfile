FROM docker.io/vhiveease/aws-python:latest
RUN pip install chameleon six futures

COPY server.py   ./
CMD ["server.handler"]
