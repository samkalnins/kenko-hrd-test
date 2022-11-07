# Kenko Heart Rate Wearable Bluetooth Grabber

Simple Go CLI to pull heartrate and interval data from bluetooth heart rate sensors.


## Compile & Run

Ensure you have golang installed from https://go.dev

Scan for peripherals first -- find your HRM UUID

    $ go run main.go -scan
    2022/11/07 14:27:33 Scanning for bluetooth devices:
    2022/11/07 14:27:43 Found 43bd4f12-09e5-303d-c21e-73270d041f7c SK WHOOP


Connect and stream data indefinitely:

    $ go run main.go -stream 43bd4f12-09e5-303d-c21e-73270d041f7c 
    2022/11/07 14:14:13 bpm=60.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:14 bpm=59.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:15 bpm=59.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:16 bpm=59.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:18 bpm=58.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:18 bpm=57.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:20 bpm=57.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:20 bpm=56.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:21 bpm=56.000000 num_rr_intervals=2 rr_intervals_ms=[1286.1328125 683.59375]
    2022/11/07 14:14:22 bpm=56.000000 num_rr_intervals=1 rr_intervals_ms=[1027.34375]
    2022/11/07 14:14:23 bpm=56.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:25 bpm=56.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:25 bpm=56.000000 num_rr_intervals=1 rr_intervals_ms=[1082.03125]
    2022/11/07 14:14:26 bpm=56.000000 num_rr_intervals=2 rr_intervals_ms=[476.5625 815.4296875]
    2022/11/07 14:14:27 bpm=56.000000 num_rr_intervals=1 rr_intervals_ms=[1063.4765625]
    2022/11/07 14:14:28 bpm=56.000000 num_rr_intervals=1 rr_intervals_ms=[1044.921875]
    2022/11/07 14:14:30 bpm=56.000000 num_rr_intervals=0 rr_intervals_ms=[]
    2022/11/07 14:14:30 bpm=55.000000 num_rr_intervals=2 rr_intervals_ms=[1000.9765625 1042.96875]


The rr_interval data I /think/ is calculated correctly, sometimes it's present or not depending on device and unknown other conditions.


## Run prebuilt binary (Mac M1)

    $ bin/kenko_bt_grabber.macos.m1 -scan
    $ bin/kenko_bt_grabber.macos.m1  -stream 43bd4f12-09e5-303d-c21e-73270d041f7c 
