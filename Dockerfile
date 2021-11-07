FROM golang:1.17 as build

WORKDIR /go/build/

COPY . .

RUN CGO_ENABLED=0 go build -o ecommerce src/*


FROM alpine

WORKDIR /hash/

ENV LISTEN_PORT 3000

COPY --from=build /go/build/ecommerce ./ecommerce
COPY --from=build /go/build/data/products.json ./data/products.json

CMD [ "./ecommerce" ]