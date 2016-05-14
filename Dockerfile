FROM nginx:latest
COPY etc/nginx.conf /etc/nginx/nginx.conf
ADD . /app/
WORKDIR /app
RUN chmod +x ./go-file-storage
RUN chmod +x ./run.sh
CMD /app/run.sh && nginx && sleep 1000
