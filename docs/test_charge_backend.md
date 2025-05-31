
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: load-test.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m0s max duration (incl. graceful stop):
              * default: 500 looping VUs for 30s (gracefulStop: 30s)


  █ THRESHOLDS 

    http_req_duration
    ✓ 'p(95)<500' p(95)=20.12ms


  █ TOTAL RESULTS 

    checks_total.......................: 15000  496.064292/s
    checks_succeeded...................: 4.09%  614 out of 15000
    checks_failed......................: 95.90% 14386 out of 15000

    ✗ status 200
      ↳  4% — ✓ 614 / ✗ 14386

    HTTP
    http_req_duration.......................................................: avg=5.58ms min=178µs med=3.73ms max=31.02ms p(90)=13.19ms p(95)=20.12ms
      { expected_response:true }............................................: avg=4.06ms min=404µs med=2.91ms max=28.24ms p(90)=8.73ms  p(95)=11.37ms
    http_req_failed.........................................................: 95.90% 14386 out of 15000
    http_reqs...............................................................: 15000  496.064292/s

    EXECUTION
    iteration_duration......................................................: avg=1s     min=1s    med=1s     max=1.06s   p(90)=1.01s   p(95)=1.02s  
    iterations..............................................................: 15000  496.064292/s
    vus.....................................................................: 500    min=500            max=500
    vus_max.................................................................: 500    min=500            max=500

    NETWORK
    data_received...........................................................: 6.4 MB 211 kB/s
    data_sent...............................................................: 1.1 MB 38 kB/s




running (0m30.2s), 000/500 VUs, 15000 complete and 0 interrupted iterations
default ✓ [ 100% ] 500 VUs  30s
