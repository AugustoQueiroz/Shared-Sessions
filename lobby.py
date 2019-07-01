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
    if request.method == "POST":

        if "user_id" not in request.form:
            result = {"sucess": False, "error": "Missing user_id parameter"}
            return json.dumps(result)

        user_id = request.form["user_id"]

        if user_id == "":
            result = {"sucess": False, "error": "Empty user_id"}
            return json.dumps(result)

        return "hi"  # TODO: return template html page

    else:
        result = {"sucess": False, "error": "invalid request. POST only"}
        return json.dumps(result)


@app.route("/joinSession/", methods=["GET", "POST"])
@app.route("/joinSession/<session_id>", methods=["GET", "POST"])
@app.route("/joinSession/<session_id>/", methods=["GET", "POST"])
def join_session(session_id=None):  # TODO: update user room with newly joined room
    if not session_id:
        result = {"sucess": False, "error": "Empty session_id"}
        return json.dumps(result)

    db = fb.database()
    rooms = db.child("rooms").order_by_key().get()
    for room in rooms.each():
        if room.key() == session_id:
            result = room.val()
            print(result)
            print(type(result))
            result["users"].append(100)  # TODO: swap to user_id later
            db.child("rooms").child(room.key()).update(result)
            return json.dumps(result)  # TODO: redirect to the room
    else:
        result = {"sucess": False, "error": "Room not found"}
        return json.dumps(result)


@app.route("/newSession/", methods=["GET", "POST"])
def new_session():  # TODO: update user room with newly created room

    db = fb.database()
    rooms = db.child("rooms").shallow().get()
    rooms_list = list(rooms.val())
    user = {"users": [75, 6543, 765]}
    while True:
        session = random.randint(1, 1000)  # behold, the might session generator
        if session not in rooms_list:
            db.child("rooms").child(session).set(user)
            break

    return "hi"  # TODO: create and redirect room


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0")
