FROM BASE_IMAGE

WORKDIR /opt/courier

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY courier /opt/courier/bin/
# COPY data/* /opt/courier/data/

EXPOSE 9001

ENTRYPOINT ["/opt/courier/bin/courier"]
CMD ["-v", "/opt/courier/db:/opt/courier/db /opt/courier/log:/opt/courier/log"]

