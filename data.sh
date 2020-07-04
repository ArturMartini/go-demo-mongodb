#!/bin/bash
echo "Execute script"
export player='{
        "name": "Cristiano Rolnado",
        "age": 35,
        "position": "st",
        "foot": "right",
        "genre": "male",
        "ranting": 4.7,
        "country": "PRT",
        "url": "https://www.youtube.com/watch?v=u5LpWA_BDSE",
        "img": "https://tmssl.akamaized.net/images/portrait/originals/8198-1568120625.jpg"
    }'

curl -XPOST -d"$player" http://localhost:8080/players