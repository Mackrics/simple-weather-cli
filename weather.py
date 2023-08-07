#! /usr/bin/python

import requests
import json
import sys
import datetime
from datetime import datetime
import pandas as pd

# Get town from command line arg----
location = sys.argv[1]


# Function to get coords
def get_lat_long(location):
    url = "https://geocoding-api.open-meteo.com/v1/search?name=" + location + \
        "&count=1&language=en&format=json"
    res = requests.get(url)
    jsonData = json.loads(res.text)
    d = jsonData['results']
    lat  = str(d[0]['latitude']) 
    long = str(d[0]['longitude'])
    coord = {"lat": lat, "long": long}
    return(coord)


# Function to get forecast
def get_forecast(lat, long):
    url = "https://api.open-meteo.com/v1/forecast?latitude=" + lat + "&longitude=" + \
            long + "&hourly=temperature_2m,rain,windspeed_10m&forecast_days=1&timezone=Europe%2FBerlin"
    response = requests.get(url)
    data = json.loads(response.text)
    forecast = data['hourly']
    df = pd.DataFrame(forecast)
    df['time'] = pd.to_datetime(df.time)
    df = df.rename(columns = {'temperature_2m': 'temperature', 
                              "windspeed_10m" : "wind"})

    df = df[df['time'] > datetime.now()]
    return(df)

# apply functions
coord = get_lat_long(location)

print(get_forecast(coord['lat'], coord['long']))
