FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/zfeed
WORKDIR /go/src/zfeed

RUN go get ./
RUN go build

CMD if [ ${APP_ENV} = production ]; \
	then \
	app; \
	else \
	go get github.com/pilu/fresh && \
	fresh; \
	fi

EXPOSE 8080