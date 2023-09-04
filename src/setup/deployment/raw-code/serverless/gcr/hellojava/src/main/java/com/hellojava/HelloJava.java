package com.hellojava;

import com.google.gson.Gson;

import java.time.Instant;

import static spark.Spark.*;

public class HelloJava {

    static String  message;

    public static void main(String args[]) {

        /*
        * Note:
        * Default port in SparkJava is 4567
        * We are using the environment variable PORT or default to 8080
        */

        port(Integer.valueOf(System.getenv().getOrDefault("PORT", "8080")));

        get("/", (req,res) -> {
            int incrementLimit = 0;
            String reqIncrementLimit = req.queryParamOrDefault("IncrementLimit", "");
            if (!reqIncrementLimit.equals("")) {
                incrementLimit = Integer.parseInt(reqIncrementLimit);
            }
            simulateWork(incrementLimit);

            return new Gson().toJson(new HelloJavaResponse("google-does-not-specify", new String[]{Long.toString(Instant.now().toEpochMilli())}));
        });
    }
    public static void simulateWork(int incrementLimit) {
         for (int i = 0; i < incrementLimit; i++) {
            Thread.onSpinWait(); // Prevent JVM/JIT optimizations from skipping the loop
        }
    }


}

