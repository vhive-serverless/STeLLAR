# Use the official maven/Java 11 image to create a build artifact.
# https://hub.docker.com/_/maven
FROM maven:3.6.3-jdk-11 as builder

# Copy local code to the container image.
WORKDIR /app
COPY pom.xml .
COPY src ./src

# Build a release artifact.
RUN mvn package -DskipTests

# Use the Official Amazon Corretto 11 image for a lean production stage of our multi-stage build.
#https://hub.docker.com/_/amazoncorretto
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM amazoncorretto:11
# Copy the jar to the production image from the builder stage.
COPY --from=builder /app/target/helloworld-0.0.1-SNAPSHOT-jar-with-dependencies.jar helloworld.jar

ENV PORT 8080

EXPOSE 8080

# Run the web service on container startup.
CMD ["java","-jar","helloworld.jar"]
