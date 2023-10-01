FROM mysql:latest

RUN chown -R mysql:root /var/lib/mysql

ENV MYSQL_DATABASE employees
ENV MYSQL_USER mike
ENV MYSQL_PASSWORD mikepass1
ENV MYSQL_ROOT_PASSWORD rootpass

ADD data.sql /etc/mysql/data.sql

RUN cp /etc/mysql/data.sql /docker-entrypoint-initdb.d

EXPOSE 3306