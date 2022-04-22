FROM BASE_IMAGE

WORKDIR /opt/authz

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY authz /opt/authz/bin/
# COPY data/* /opt/authz/data/

EXPOSE 9001

ENTRYPOINT ["/opt/authz/bin/authz"]
CMD ["-v", "/opt/authz/db:/opt/authz/db /opt/authz/log:/opt/authz/log"]

