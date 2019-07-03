FROM python:3.7

COPY requirements.txt ./

RUN pip install -r requirements.txt

COPY . .

EXPOSE 3000

CMD ["python3", "lobby.py"]