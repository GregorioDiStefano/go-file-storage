FROM nginx:latest
RUN apt-get update && apt-get install -y supervisor
ENV CONFIG_FILE ./config/config.testing.json
COPY etc/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
COPY etc/nginx.conf /etc/nginx/nginx.conf
ADD * /app/

WORKDIR /app
RUN chmod +x ./http-file-sharing

CMD ["/usr/bin/supervisord"]
