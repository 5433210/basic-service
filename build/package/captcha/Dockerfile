FROM BASE_IMAGE

WORKDIR /opt/captcha

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
      echo "Asia/Shanghai" > /etc/timezone

COPY captcha /opt/captcha/
COPY captcha.yaml /opt/captcha/

EXPOSE 3003

ENTRYPOINT ["/opt/captcha/captcha"]

