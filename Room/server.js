var firebase = require("firebase");
var SpotifyWebApi = require('spotify-web-api-node');

var roomToken = process.env.token_room || '1';

var firebaseConfig = {
  apiKey: "AIzaSyAowSk2cCFSnfBCqjUOW_SPOokjjNjvyoo",
  authDomain: "papo-5edd4.firebaseapp.com",
  databaseURL: "https://papo-5edd4.firebaseio.com",
  projectId: "papo-5edd4",
  storageBucket: "papo-5edd4.appspot.com",
  messagingSenderId: "820296196555",
  appId: "1:820296196555:web:d9d8cffd9df33225"
};

var defaultProj = firebase.initializeApp(firebaseConfig);

var database = firebase.database();

var payload = {};

var intervalId = setInterval(function () {
  database.ref('rooms/' + roomToken).once('value').then(function (snap) {
    tokens = snap.val().users;

    database.ref('users/' + tokens[0]).once('value').then(function (volta) {
      //console.log(volta.val());
      token = volta.val().token.access_token;




      var spotifyApi = new SpotifyWebApi({
        redirectUri: '0.0.0.0:8000',
        clientId: '32cca9f1a35c44d4bed142d7fe78a3a8'
      })
      //console.log(tokens[0]);
      spotifyApi.setAccessToken(token);
      spotifyApi.getMyCurrentPlaybackState({
      })
        .then(function (data) {
          // Output items
          console.log(data.body.item);
          console.log(data.body.progress_ms);

          console.log("Now Playing: ", data.body);
          payload = {
            "context_uri": data.body.item.uri,
            "offset": {
              "position": data.body.item.track_number
            },
            "position_ms": data.body.progress_ms
          };
        }, function (err) {
          console.log('Something went wrong!', err);
        });
    }, function (err) { console.log(err) });



    for (i = 1; i < tokens.length; i++) {
       database.ref('users/' + tokens[i]).once('value').then(function (volta) {
      token = volta.val().token.access_token;
      spotifyApi.setAccessToken(token);
      spotifyApi.getMyCurrentPlaybackState({
      })
        .then(function (data) {
          // Output items
          console.log(data.body.item);
          console.log(data.body.progress_ms);

          console.log("Now Playing: ", data.body);
          /*payload = {
            "context_uri": "spotify:album:5ht7ItJgpBH7W6vJ5BqpPr",
            "offset": {
              "position": 6
            },
            "position_ms": 20000
          };*/
          if (data.body.item.uri != payload["context_uri"]) {


            spotifyApi.play(payload)//({context_uri:'spotify:track:0CZ8lquoTX2Dkg7Ak2inwA',offset:5, position_ms:'199978'})
              .then(function (yay) {
                //whatever
                console.log(yay);
              }, function (err) {
                console.log(err);
              });
          }
        }, function (err) {
          console.log('Something went wrong!', err);
        });
      }, function(err){console.log(err)});
    }
  



  });
}, 5000);