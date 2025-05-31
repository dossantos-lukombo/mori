import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 500,
  duration: '30s',
  thresholds: {
    'http_req_duration': ['p(95)<500'],
  },
};

// Variante A : fonction anonyme exportée par défaut
export default function () {
  const res = http.get('http://localhost:8081/health');
  check(res, { 'status 200': (r) => r.status === 200 });
  sleep(1);
}
