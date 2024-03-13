package com.hellojava;

import spark.Request;

public class HelloJavaResponse {
    private String RequestID;
    private String[] TimestampChain;

    public HelloJavaResponse(String RequestID, String[] TimeStampChain){
        this.RequestID = RequestID;
        this.TimestampChain = TimeStampChain;
    }
}
