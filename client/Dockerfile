FROM erlang:25

WORKDIR /tmp
RUN apt-get update && \
    apt-get install -y python perl libtemplate-perl gnuplot && \
    wget http://tsung.erlang-projects.org/dist/tsung-1.7.0.tar.gz

RUN tar xzf tsung-1.7.0.tar.gz

WORKDIR /tmp/tsung-1.7.0
RUN ./configure && \
    make && \
    make install

RUN mkdir -p /usr/src/app
RUN cp ./tsung-1.0.dtd /usr/src/app

WORKDIR /tmp
RUN rm -rf ./tsung*

WORKDIR /usr/src/app