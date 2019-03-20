FROM alpine:3.4

RUN apk -U add ca-certificates ffmpeg

RUN mkdir /var/www && mkdir /var/www/.cache

EXPOSE 1333
EXPOSE 443
EXPOSE 80

ADD media-check-tool /bin/media-check-tool

CMD ["media-check-tool"]