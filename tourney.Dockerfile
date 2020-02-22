FROM golang:latest  as go-builder

RUN mkdir /src && mkdir /output

ADD . /src

WORKDIR /src

RUN go build -o /output/get5-wrapper ./cmd/build-config/

FROM lacledeslan/gamesvr-csgo-tourney:get5

COPY --chown=CSGOTourneyGet5:root --from=go-builder /output /app

COPY --chown=CSGOTourneyGet5:root --from=go-builder /src/cmd/build-config/get5-wrapper.json /app/get5-wrapper.json
