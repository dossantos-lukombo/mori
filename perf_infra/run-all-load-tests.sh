#!/usr/bin/env bash
set -e

# 1) Créer un dossier "results" s'il n'existe pas déjà,
#    pour y déposer uniquement les fichiers de résumé de tests.
mkdir -p results

# 2) Déclarer deux tableaux parallèles : noms de services et leurs endpoints
SERVICES=("go-backend" "frontend" "llm-server")
ENDPOINTS=(
  "http://localhost:8081/health"
  "http://localhost:8080/"
  "http://localhost:3000/health"
)

# Vérification que SERVICES et ENDPOINTS ont la même longueur
if [ "${#SERVICES[@]}" -ne "${#ENDPOINTS[@]}" ]; then
  echo "Erreur : SERVICES et ENDPOINTS doivent avoir le même nombre d'éléments." >&2
  exit 1
fi

# 3) Liste des paliers de VUs à tester
VUS_LIST=(500 1000 1500 2000 2500 3000)

# 4) Pour chaque service ET chaque palier de VUs,
#    on exécute k6 en mode "summary-export" pour n’avoir que le résultat final
for idx in "${!SERVICES[@]}"; do
  SERVICE="${SERVICES[$idx]}"
  URL="${ENDPOINTS[$idx]}"

  for VUS in "${VUS_LIST[@]}"; do
    SUMMARY_FILE="results/${SERVICE}-${VUS}-summary.json"

    echo
    echo "────────────────────────────────────────────────────────"
    echo "▶ Test de charge (séquentiel) : ${SERVICE} – ${VUS} VUs – 30s"
    echo "  → Résumé exporté dans : ${SUMMARY_FILE}"
    echo "────────────────────────────────────────────────────────"
    echo

    # On lance k6 en utilisant uniquement --summary-export
    # Cela génère un fichier JSON contenant le résumé (p(50), p(95), RPS, taux d'erreur…) 
    # pour ce service à ce niveau de charge.
    TARGET_URL="$URL" VUS="$VUS" \
      k6 run single-service-test.js --summary-export "${SUMMARY_FILE}"

    echo
    echo ">>> Fin du test ${SERVICE} (${VUS} VUs). Résumé dans ${SUMMARY_FILE}."
    echo

    # Pause courte facultative pour laisser retomber la charge
    sleep 3
  done
done

echo "✅ Tous les tests séquentiels sont terminés. Résumés dans le dossier 'results/'."
