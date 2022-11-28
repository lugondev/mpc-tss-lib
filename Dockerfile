# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/* \

WORKDIR /app
# Copy the binary to the production image from the builder stage.
COPY ./configuration.yml ./configuration.yml
COPY ./build/server ./server
ENV MULTISIG_DB_PROFILE postgresql.primary

#RUN ls -alh .
# Run the web service on container startup.

CMD ["./server"]

# [END run_helloworld_dockerfile]
# [END cloudrun_helloworld_dockerfile]
