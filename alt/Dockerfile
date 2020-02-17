FROM python:3.8.1-slim-buster

WORKDIR /app/src

USER root
RUN apt-get update && apt-get upgrade -y

RUN groupadd  -g 8888 app
RUN useradd -u 8888 -d /app/src -s /bin/bash -g app -M app
RUN mkdir -p /app/src && chown -R app:app /app

RUN pip install --no-cache-dir pip-tools python-dotenv

USER app
ENV HOME=/app/src

RUN install -d -o app -g app /app && install -d -o app -g app /app/src
COPY --chown=app:app requirements.txt /app/src
RUN python3 -m venv /app/venv && . /app/venv/bin/activate && pip install --no-cache-dir -r /app/src/requirements.txt

USER root
RUN . /app/venv/bin/activate && https -F -I -o /tini github.com/krallin/tini/releases/download/v0.18.0/tini-static-amd64 && chmod 755 /tini

USER app

COPY --chown=app:app src /app/src

COPY service /service

ENTRYPOINT [ "/tini", "--", "/service" ]
