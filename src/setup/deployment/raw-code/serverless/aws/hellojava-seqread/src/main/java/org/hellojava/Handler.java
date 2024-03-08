package org.hellojava;

import java.util.Map;
import java.util.HashMap;
import java.time.Instant;
import java.util.Random;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyRequestEvent;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyResponseEvent;
import com.google.gson.Gson;
import org.crac.Resource;
import org.crac.Core;

class ResponseEventBody {
    String region;
    String requestId;
    String[] timestampChain;

    public ResponseEventBody(String region, String requestId, String[] timestampChain) {
        this.region = region;
        this.requestId = requestId;
        this.timestampChain = timestampChain;
    }
}

public class Handler implements RequestHandler<APIGatewayProxyRequestEvent, APIGatewayProxyResponseEvent>, Resource{
    byte[] pageData;

    public Handler() {
        Core.getGlobalContext().register(this);
    }
    @Override
    public void beforeCheckpoint(org.crac.Context<? extends Resource> context) throws Exception {
        int size = (int)(1024 * 1024 * 400);
        this.pageData = new byte[size];
        new Random().nextBytes(this.pageData);
    }

    @Override
    public void afterRestore(org.crac.Context<? extends Resource> context) {

    }

    @Override
    public APIGatewayProxyResponseEvent handleRequest(APIGatewayProxyRequestEvent event, Context context)
    {
        Gson gson = new Gson();
        int incrementLimit = 0;
        if (event.getQueryStringParameters() != null) {
            incrementLimit = Integer.parseInt(event.getQueryStringParameters().getOrDefault("incrementLimit", "0"));
        }
        this.simulateWork(incrementLimit);
        String requestId = "no-context";
        if (context != null) {
            requestId = context.getAwsRequestId();
        }

        long startTime = System.currentTimeMillis();
        readDataSeq();
        long endTime = System.currentTimeMillis();

        long pageReadTime = endTime-startTime;

        String[] timestampChain = new String[]{""+pageReadTime};
        ResponseEventBody resBody = new ResponseEventBody(System.getenv("AWS_REGION"), requestId, timestampChain);

        Map<String, String> responseHeaders = new HashMap<>();
        responseHeaders.put("Content-Type", "application/json");
        APIGatewayProxyResponseEvent response = new APIGatewayProxyResponseEvent().withHeaders(responseHeaders);
        response.setIsBase64Encoded(false);
        response.setStatusCode(200);
        response.setBody(gson.toJson(resBody));

        return response;
    }

    public void simulateWork(int incrementLimit) {
        int i = 0;
        while (i < incrementLimit) {
            i++;
        }
    }

    public void readDataSeq() {
        for (int i = 0; i< this.pageData.length; i+=4096) {
            System.out.print(this.pageData[i]);
        }
    }
}