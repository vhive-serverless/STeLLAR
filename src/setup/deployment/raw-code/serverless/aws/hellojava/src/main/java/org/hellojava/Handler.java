package org.hellojava;

import java.util.Map;
import java.util.HashMap;
import java.time.Instant;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyRequestEvent;
import com.amazonaws.services.lambda.runtime.events.APIGatewayProxyResponseEvent;
import com.google.gson.Gson;

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

public class Handler implements RequestHandler<APIGatewayProxyRequestEvent, APIGatewayProxyResponseEvent>{

    @Override
    public APIGatewayProxyResponseEvent handleRequest(APIGatewayProxyRequestEvent event, Context context)
    {
        Gson gson = new Gson();
        int incrementLimit = Integer.parseInt(event.getQueryStringParameters().getOrDefault("incrementLimit", "0"));
        this.simulateWork(incrementLimit);
        String requestId = "no-context";
        if (context != null) {
            requestId = context.getAwsRequestId();
        }

        Instant now = Instant.now();
        String[] timestampChain = new String[]{""+now.getEpochSecond()+now.getNano()};
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
}
