import timeit

import pandas as pd
from sklearn.model_selection import train_test_split
from sklearn.ensemble import RandomForestRegressor

# data cleansing and manipulation
def predict_flight_delays():
    flights = pd.read_csv('./flights.csv',)
    airports = pd.read_csv('./airports.csv')

    flights.isnull().values.any()
    flights.isnull().sum()

    variables_to_remove=["YEAR","FLIGHT_NUMBER","TAIL_NUMBER","DEPARTURE_TIME","TAXI_OUT","WHEELS_OFF","ELAPSED_TIME","AIR_TIME","WHEELS_ON","TAXI_IN","ARRIVAL_TIME","DIVERTED","CANCELLED","CANCELLATION_REASON","AIR_SYSTEM_DELAY", "SECURITY_DELAY","AIRLINE_DELAY","LATE_AIRCRAFT_DELAY","WEATHER_DELAY","SCHEDULED_TIME","SCHEDULED_ARRIVAL"]

    flights.drop(variables_to_remove, axis=1, inplace= True)

    flights.columns

    flights.loc[~flights.ORIGIN_AIRPORT.isin(airports.IATA_CODE.values),'ORIGIN_AIRPORT']='OTHER'
    flights.loc[~flights.DESTINATION_AIRPORT.isin(airports.IATA_CODE.values),'DESTINATION_AIRPORT']='OTHER'

    flights=flights.dropna()

    flights['DAY_OF_WEEK'] = flights['DAY_OF_WEEK'].apply(str)
    flights['DAY_OF_WEEK'].replace({"7":"SUNDAY", "1": "MONDAY", "2": "TUESDAY", "3":"WEDNESDAY", "4":"THURSDAY", "5":"FRIDAY", "6":"SATURDAY"}, inplace=True)

    dums = ['AIRLINE','ORIGIN_AIRPORT','DESTINATION_AIRPORT','DAY_OF_WEEK']
    flights_cat=pd.get_dummies(flights[dums],drop_first=True)

    var_to_remove=["DAY_OF_WEEK","AIRLINE","ORIGIN_AIRPORT","DESTINATION_AIRPORT"]
    flights.drop(var_to_remove, axis=1, inplace=True)

    cleansedData = pd.concat([flights, flights_cat],axis=1)
    print(cleansedData.columns)

    # model training
    X=cleansedData.drop("DEPARTURE_DELAY",axis=1)
    Y=cleansedData.DEPARTURE_DELAY

    X_train, X_test, Y_train, Y_test = train_test_split(X, Y, test_size=0.3, random_state=0)

    reg_rf = RandomForestRegressor()
    reg_rf.fit(X_train,Y_train)

    print(reg_rf.score(X_test,Y_test))

print(f"execution duration {timeit.timeit(predict_flight_delays, number=1)}")