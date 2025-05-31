# Documentation du Processus CI/CD “Développer → Production”

> _« On ne pousse jamais directement sur `develop` ni sur `production`. Tout passe par PR, même la magie du bump de version. »_

---

## Table des matières

1. [Contexte général](#contexte-général)  
2. [Stratégie de branches et protections](#stratégie-de-branches-et-protections)  
3. [Bump de version automatique sur PR vers `develop`](#bump-de-version-automatique-sur-pr-vers-develop)  
   1. [Principe et prérequis](#principe-et-prérequis)  
   2. [Configuration de semantic-release](#configuration-de-semantic-release)  
   3. [Workflow GitHub Actions pour bump-auto](#workflow-github-actions-pour-bump-auto)  
4. [CI (Tests) sur PR vers `develop`](#ci-tests-sur-pr-vers-develop)  
   1. [Objectif](#objectif)  
   2. [Workflow GitHub Actions – ci-test.yml](#workflow-github-actions--ci-testyml)  
5. [CD (Build & Déploiement) via PR `develop → production`](#cd-build--déploiement-via-pr-develop--production)  
   1. [Objectif](#objectif-1)  
   2. [Workflow GitHub Actions – cd-deploy.yml](#workflow-github-actions--cd-deployyml)  
6. [Fichier `docker-compose.yml` utilisé en production](#fichier-docker-composeyml-utilisé-en-production)  
7. [Résumé du flux complet](#résumé-du-flux-complet)  
8. [Quelques conseils pour rouler tranquille](#quelques-conseils-pour-rouler-tranquille)  

---

## 1. Contexte général

Vous avez trois services Docker qui tournent dans un seul dépôt :

- **go-backend** (sur le port 8081)  
- **llm-server** (sur le port 3000)  
- **frontend** (sur le port 8080)  

Vous voulez :

1. **Ne jamais pousser directement** sur `develop` ni sur `production` (branches protégées).  
2. Bumper la version **automatiquement** dès qu’une PR vers `develop` est mergée.  
3. Déclencher le déploiement en production **uniquement** quand une PR fusionne `develop` → `production`.  
4. Tester en CI le bon fonctionnement de chaque container avant tout merge.  
5. Garder un historique clair des versions (`vX.Y.Z`) et des releases GitHub.  

> Pour faire simple : on MARGE via PR, on laisse GitHub Actions faire le reste.  

---

## 2. Stratégie de branches et protections

1. **Branche `develop`**  
   - Protégée : on ne peut pas y pousser directement.  
   - Seule façon d’y faire des modifications : PR validée.  
   - Quand une PR vers `develop` est **fermée** (event `pull_request.closed`) **ET** mergée, on déclenche le **bump de version automatique**.  

2. **Branche `production`**  
   - Protégée : on ne peut pas y pousser directement.  
   - Seule façon d’y faire des modifications : PR validée venant de `develop`.  
   - Quand une PR `develop` → `production` est **fermée** (event `pull_request.closed`) **ET** mergée, on déclenche le **CD** (build Docker → push sur Docker Hub → déploiement sur le VPS).  

3. **Tags Git**  
   - À chaque PR mergée sur `develop`, l’action de bump crée un commit de type `chore(release): vX.Y.Z` et pose un **tag Git** `vX.Y.Z` sur la branche `develop`.  
   - Ce tag devient disponible pour le workflow de CD.  

> **Avertissement** : le repo doit accorder à GitHub Actions les droits d’écriture sur les branches protégées pour que `semantic-release` puisse pousser le commit de bump et tagger.

---

## 3. Bump de version automatique sur PR vers `develop`

### 3.1 Principe et prérequis

1. **Commits “Conventional Commits”**  
   - **feat(scope): description** → bump MINOR ( X.Y.Z → X.(Y+1).0 ).  
   - **fix(scope): description** → bump PATCH ( X.Y.Z → X.Y.(Z+1) ).  
   - **BREAKING CHANGE** (dans le corps du commit) → bump MAJOR ( X.Y.Z → (X+1).0.0 ).  

2. **semantic-release**  
   - Outil Node.js qui analyse l’historique des commits depuis le dernier tag, détermine automatiquement la nouvelle version sémantique (MAJOR/MINOR/PATCH), met à jour le `CHANGELOG.md` et le fichier de version (`package.json` ou `VERSION`), crée le commit `chore(release): vX.Y.Z` et le tag `vX.Y.Z`, puis pousse tout sur la branche.  
   - **Exige** que la branche soit accessible en écriture par l’action (options “Allow GitHub Actions to push” activée dans les règles de protection).  

3. **GitHub Action “pull_request.closed”**  
   - On cible l’événement `pull_request.closed` sur `develop`. Si `merged == true`, on lance `semantic-release`.  
   - Le commit de bump (et le tag) sont créés **après** la fusion de la PR, de sorte que le merge commit inclut les dernières fonctionnalités.  

### 3.2 Configuration de semantic-release

1. **Installer les dépendances** (à la racine du repo, ou dans un dossier `backend`, `frontend`, selon la stack) :  
   ```bash
   npm install --save-dev semantic-release \
     @semantic-release/changelog \
     @semantic-release/git \
     @semantic-release/github \
     @semantic-release/commit-analyzer \
     @semantic-release/release-notes-generator
