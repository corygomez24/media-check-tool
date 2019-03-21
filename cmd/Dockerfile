FROM alpine:3.4

RUN apk -U add ca-certificates ffmpeg

RUN mkdir /var/www && mkdir /var/www/.cache

EXPOSE 1333
EXPOSE 443
EXPOSE 80

ADD media-check-tool /bin/media-check-tool
RUN mkdir pkg && mkdir pkg/media-check-tool
COPY AudioOnly.png /pkg/media-check-tool/AudioOnly.png

CMD ["media-check-tool"]