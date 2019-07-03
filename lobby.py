from flask import Flask, request, render_template
import json
import os
import random

import spotipy
import pyrebase
from dotenv import load_dotenv


load_dotenv()

app = Flask(__name__)

fb_config = {
    "apiKey": os.getenv("FB_API_KEY"),
    "authDomain": os.getenv("FB_AUTH_DOMAIN"),
    "databaseURL": os.getenv("FB_DB_URL"),
    "storageBucket": "",
}

fb = pyrebase.initialize_app(fb_config)

sp = spotipy.Spotify()


@app.route("/", methods=["GET", "POST"])
def user_in():

    user_id = request.args.get("user_id", "")
    if user_id == "":
        result = {"sucess": False, "error": "Empty user_id"}
        return json.dumps(result)

    return render_template("index.html", user=user_id)


@app.route("/joinSession/", methods=["GET", "POST"])
def join_session():  # TODO: update user room with newly joined room

    user_id = request.args.get("user_id", "")
    if user_id == "":
        result = {"sucess": False, "error": "Empty user_id"}
        return json.dumps(result)

    session_id = request.args.get("sessionCode", "")
    if session_id == "":
        result = {"sucess": False, "error": "Empty session_id"}
        return json.dumps(result)

    db = fb.database()
    rooms = db.child("rooms").order_by_key().get()
    for room in rooms.each():
        if room.key() == session_id:
            result = room.val()
            result["users"].append(100)  # TODO: swap to user_id later
            db.child("rooms").child(room.key()).update(result)
            return json.dumps(result)  # TODO: redirect to the room
    else:
        result = {"sucess": False, "error": "Room not found"}
        return json.dumps(result)


@app.route("/newSession/", methods=["GET", "POST"])
def new_session():  # TODO: update user room with newly created room

    user_id = request.args.get("user_id", "")
    if user_id == "":
        result = {"sucess": False, "error": "Empty user_id"}
        return json.dumps(result)

    db = fb.database()
    rooms = db.child("rooms").shallow().get()
    rooms_list = list(rooms.val())
    user = {"users": [int(user_id)]}
    session = 0
    while True:
        session = random.randint(1, 1000)  # behold, the might session generator
        if session not in rooms_list:
            db.child("rooms").child(session).set(user)
            break

    db.child("users").child(int(user_id)).update({"room": session})

    return "hi"  # TODO: create and redirect room


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0")
