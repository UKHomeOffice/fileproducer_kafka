FROM scratch
MAINTAINER Ivan Pedrazas <ipedrazas@gmail.com>

ADD ./file_producer /file_producer

ENTRYPOINT ["/file_producer"]


