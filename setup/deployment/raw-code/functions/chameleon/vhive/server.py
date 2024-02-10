import grpc
import logging
import six
from chameleon import PageTemplate
from concurrent import futures

import chainfunction_pb2
import chainfunction_pb2_grpc

BIGTABLE_ZPT = """\
<table xmlns="http://www.w3.org/1999/xhtml"
xmlns:tal="http://xml.zope.org/namespaces/tal">
<tr tal:repeat="row python: options['table']">
<td tal:repeat="c python: row.values()">
<span tal:define="d python: c + 1"
tal:attributes="class python: 'column-' + %s(d)"
tal:content="python: d" />
</td>
</tr>
</table>""" % six.text_type.__name__

responses = ["record_response", "replay_response"]


class Greeter(chainfunction_pb2_grpc.ProducerConsumerServicer):

    def InvokeNext(self, _, context):
        tmpl = PageTemplate(BIGTABLE_ZPT)

        data = {}
        num_of_cols = 15
        num_of_rows = 10

        for i in range(num_of_cols):
            data[str(i)] = i

        table = [data for x in range(num_of_rows)]
        options = {'table': table}

        data = tmpl.render(options=options)

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
