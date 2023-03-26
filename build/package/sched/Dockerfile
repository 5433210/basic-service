FROM BASE_IMAGE

WORKDIR /opt/sched

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY sched /opt/sched/
COPY sched.yaml /opt/sched/
COPY mqtt_sender /opt/sched/

EXPOSE 3005

ENTRYPOINT ["/opt/sched/sched"]

