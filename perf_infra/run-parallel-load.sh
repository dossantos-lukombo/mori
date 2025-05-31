#!/usr/bin/env bash
set -e

# 1) Créer un dossier "results" s'il n'existe pas déjà,
#    pour y déposer le résumé global.
mkdir -p results

# 2) Nom du fichier de résumé global pour tous les scénarios
SUMMARY_FILE="results/parallel-load-test-summary.json"

echo
echo "────────────────────────────────────────────────────────"
echo "▶ Lancement du test parallèle k6 pour tous les services"
echo "  → Résumé global exporté dans : ${SUMMARY_FILE}"
echo "────────────────────────────────────────────────────────"
echo

# 3) Exécuter k6 avec --summary-export :
#    cela lancera simultanément les 9 scénarios définis dans load-test-parallel.js
#    et n’écrira qu’un seul JSON résumé.
k6 run load-test-parallel.js --summary-export "${SUMMARY_FILE}"

echo
echo "✅ Test parallèle terminé. Fichier résumé disponible : ${SUMMARY_FILE}"
echo
