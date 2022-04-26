FROM BASE_IMAGE

WORKDIR /opt/authn

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY authn /opt/authn/
COPY authn.yaml /opt/authn/

EXPOSE 3001

ENTRYPOINT ["/opt/authn/authn"]

