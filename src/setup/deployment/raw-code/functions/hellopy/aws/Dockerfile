FROM docker.io/vhiveease/aws-python:latest
RUN pip install futures

COPY server.py   ./
CMD ["server.handler"]
