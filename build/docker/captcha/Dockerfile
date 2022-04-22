FROM BASE_IMAGE

WORKDIR /opt/captcha

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY captcha /opt/captcha/bin/
# COPY data/* /opt/captcha/data/

EXPOSE 9001

ENTRYPOINT ["/opt/captcha/bin/captcha"]
CMD ["-v", "/opt/captcha/db:/opt/captcha/db /opt/captcha/log:/opt/captcha/log"]

