import http from 'k6/http';
import { check, sleep } from 'k6';

// k6 met automatiquement la variable d'env __ENV.TARGET_URL et __ENV.VUS
export let options = {
  vus: __ENV.VUS ? parseInt(__ENV.VUS) : 100, // 100 VU par défaut
  duration: '30s',
  thresholds: {
    'http_req_duration': ['p(95)<500'],
  },
};

export default function () {
  const res = http.get(__ENV.TARGET_URL);
  check(res, { 'status 200': (r) => r.status === 200 });
  sleep(1);
}
