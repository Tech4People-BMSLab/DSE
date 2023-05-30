FROM python:3.8

COPY ./api/requirements.txt /api/requirements.txt
COPY ./api/data             /api/data/
COPY ./api/main.py          /api/main.py

RUN python3 -m pip install --upgrade pip
RUN python3 -m pip install --upgrade setuptools
# RUN python3 -m pip install -r requirements.txt

CMD [ "python", "/api/main.py" ]