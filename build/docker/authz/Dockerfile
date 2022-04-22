FROM BASE_IMAGE

WORKDIR /opt/authz

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY authz /opt/authz/
COPY authz.yaml /opt/authz/
COPY authz.rego /opt/authz/

EXPOSE 3002

ENTRYPOINT ["/opt/authz/authz"]

