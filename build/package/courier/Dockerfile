FROM BASE_IMAGE

WORKDIR /opt/courier

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY courier /opt/courier/
COPY courier.yaml /opt/courier/

EXPOSE 3004

ENTRYPOINT ["/opt/courier/courier"]

