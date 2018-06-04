FROM bitriseio/docker-bitrise-base

# envs
ENV PROJ_NAME=bitrise-step-analytics
ENV BITRISE_SOURCE_DIR="/bitrise/go/src/github.com/slapec93/$PROJ_NAME"

# Get go tools
RUN go get github.com/codegangsta/gin \
    && go get github.com/kisielk/errcheck \
    && go get github.com/golang/lint/golint \
    && go get github.com/stripe/safesql

WORKDIR $BITRISE_SOURCE_DIR

CMD $PROJ_NAME
