FROM letfn/python

COPY --chown=app:app requirements.txt /app/src/
RUN . /app/venv/bin/activate && pip install --no-cache-dir -r /app/src/requirements.txt

COPY service /service

ENTRYPOINT [ "/tini", "--", "/service" ]
