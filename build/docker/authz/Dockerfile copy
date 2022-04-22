FROM BASE_IMAGE

WORKDIR /opt/authn

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY authn /opt/authn/bin/
# COPY data/* /opt/authn/data/

EXPOSE 9001

ENTRYPOINT ["/opt/authn/bin/authn"]
CMD ["-v", "/opt/authn/db:/opt/authn/db /opt/authn/log:/opt/authn/log"]

