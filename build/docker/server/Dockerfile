FROM BASE_IMAGE

WORKDIR /opt/rbac

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY server /opt/rbac/bin/
COPY data/* /opt/rbac/data/

EXPOSE 9001

ENTRYPOINT ["/opt/rbac/bin/server"]
CMD ["-v", "/opt/rbac/db:/opt/rbac/db /opt/rbac/log:/opt/rbac/log"]

