FROM body-controller-monorepo

WORKDIR /domains/security/services/aggregation/

RUN go build -mod=mod -o service ./*.go

CMD [ "./service" ]
