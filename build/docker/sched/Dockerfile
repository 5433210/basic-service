FROM BASE_IMAGE

WORKDIR /opt/sched

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY sched /opt/sched/bin/
# COPY data/* /opt/sched/data/

EXPOSE 9001

ENTRYPOINT ["/opt/sched/bin/sched"]
CMD ["-v", "/opt/sched/db:/opt/sched/db /opt/sched/log:/opt/sched/log"]

