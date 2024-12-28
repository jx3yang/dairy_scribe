# Diary Scribe

## Run the Server
```
cd events_aggregator/
go build
./events_aggregator
```

## Run the MacOS App
```
cd events_logger/macos_app/
make
./events_logger_app
```

## Run the Browser Extension
```
cd events_logger/chrome_extension/
yarn && yarn build
```
Import the extension at `chrome://extensions` by selecting the `events_logger/chrome_extension/dist/` directory

## Run the Scribe
```
cd scribe/
go build
./scribe <path_to_logs>
```

## Design
![image](https://github.com/user-attachments/assets/f95d85b8-bfb1-4d1a-a071-e6017cdfc0c5)

