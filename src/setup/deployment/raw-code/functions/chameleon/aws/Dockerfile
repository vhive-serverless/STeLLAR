FROM docker.io/vhiveease/aws-python:latest
RUN pip install chameleon

COPY server.py   ./
CMD ["server.handler"]
