import http from 'k6/http';
import { check, sleep } from 'k6';

/*
  9 scénarios parallèles :
    - go-backend : 500, 1000, 1500 VUs
    - frontend   : 500, 1000, 1500 VUs
    - llm-server : 500, 1000, 1500 VUs

  Chaque scénario tourne 30 secondes en “constant-vus” (nombre fixe de VUs).
  On effectue un simple GET sur un endpoint “health” ou “/” et on vérifie le code 200.
  On met un sleep(1) entre chaque itération pour simuler un utilisateur qui attend 1 s.
*/

export let options = {
  scenarios: {
    // ─── go-backend ─────────────────────────────────────────────────────
    go_500: {
      executor: 'constant-vus',
      exec: 'testGo',
      vus: 500,
      duration: '30s',
    },
    go_1000: {
      executor: 'constant-vus',
      exec: 'testGo',
      vus: 1000,
      duration: '30s',
    },
    go_1500: {
      executor: 'constant-vus',
      exec: 'testGo',
      vus: 1500,
      duration: '30s',
    },

    // ─── frontend ───────────────────────────────────────────────────────
    fe_500: {
      executor: 'constant-vus',
      exec: 'testFront',
      vus: 500,
      duration: '30s',
    },
    fe_1000: {
      executor: 'constant-vus',
      exec: 'testFront',
      vus: 1000,
      duration: '30s',
    },
    fe_1500: {
      executor: 'constant-vus',
      exec: 'testFront',
      vus: 1500,
      duration: '30s',
    },

    // ─── llm-server ─────────────────────────────────────────────────────
    llm_500: {
      executor: 'constant-vus',
      exec: 'testLlm',
      vus: 500,
      duration: '30s',
    },
    llm_1000: {
      executor: 'constant-vus',
      exec: 'testLlm',
      vus: 1000,
      duration: '30s',
    },
    llm_1500: {
      executor: 'constant-vus',
      exec: 'testLlm',
      vus: 1500,
      duration: '30s',
    },
  },
  thresholds: {
    // Exemple de seuils : 95% des requêtes < 500 ms pour chaque scénario
    'http_req_duration{scenario:go_500}':    ['p(95)<500'],
    'http_req_duration{scenario:go_1000}':   ['p(95)<500'],
    'http_req_duration{scenario:go_1500}':   ['p(95)<500'],
    'http_req_duration{scenario:fe_500}':    ['p(95)<500'],
    'http_req_duration{scenario:fe_1000}':   ['p(95)<500'],
    'http_req_duration{scenario:fe_1500}':   ['p(95)<500'],
    'http_req_duration{scenario:llm_500}':   ['p(95)<500'],
    'http_req_duration{scenario:llm_1000}':  ['p(95)<500'],
    'http_req_duration{scenario:llm_1500}':  ['p(95)<500'],
  },
};

// ─── Fonctions d’exécution pour chaque service ─────────────────────────────────

export function testGo() {
  // Remplacez par l’URL “health” ou "/" de votre go-backend si besoin
  const res = http.get('http://localhost:8081/health');
  check(res, { 'go-backend status 200': (r) => r.status === 200 });
  sleep(1);
}

export function testFront() {
  // Remplacez par l’URL racine de votre frontend (par ex. page d’accueil)
  const res = http.get('http://localhost:8080/');
  check(res, { 'frontend status 200': (r) => r.status === 200 });
  sleep(1);
}

export function testLlm() {
  // Remplacez par l’URL “health” de votre llm-server (ou "/" si pas de health)
  const res = http.get('http://localhost:3000/health');
  check(res, { 'llm-server status 200': (r) => r.status === 200 });
  sleep(1);
}
