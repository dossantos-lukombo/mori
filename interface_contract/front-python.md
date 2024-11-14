# Contrat d'Interface pour l'API `/api/process-message`

## Description
Ce document décrit les structures de données utilisées pour envoyer et recevoir des requêtes sur l'endpoint `/api/process-message` de l'API. Cet endpoint permet de traiter les messages des utilisateurs et de renvoyer une réponse appropriée.

## Requête Entrante

### Classe `Data`
La classe `Data` représente la structure des données attendues par l'API lorsqu'une requête est envoyée par le client.

- **Type** : `application/json`
- **Structure** :

```json
{
  "user_id": "string",
  "conversation_id": "string",
  "message": "string"
}
```

### Description des Champs

- **`user_id` (str)** : Identifiant unique de l'utilisateur qui envoie la requête.
- **`conversation_id` (str)** : Identifiant unique de la conversation pour contextualiser le message.
- **`message` (str)** : Le message ou la requête de l'utilisateur.

## Réponse Sortante

### Classe `SendData`
La classe `SendData` représente la structure des données renvoyées par l'API après traitement de la requête.

- **Type** : `application/json`
- **Structure** :

```json
{
  "status": "success" | "error",
  "user_id": "string",
  "conversation_id": "string",
  "response": "string",
  "timestamp": "string" | format: "2024-11-14T14:32:00Z"
}
```

### Description des Champs

- **`status` (Literal["success", "error"])** : Statut de la réponse. Peut être `success` pour indiquer un traitement réussi ou `error` pour indiquer un échec.
- **`user_id` (str)** : Identifiant unique de l'utilisateur correspondant à la requête initiale.
- **`conversation_id` (str)** : Identifiant unique de la conversation, identique à celui de la requête.
- **`response` (str)** : Contenu de la réponse générée par le système.
- **`timestamp` (datetime)** : Horodatage de la réponse au format ISO 8601.

## Exemple de Requête

```json
POST /api/process-message
Content-Type: application/json

{
  "user_id": "12345",
  "conversation_id": "67890",
  "message": "Bonjour, pouvez-vous m'aider ?"
}
```

## Exemple de Réponse

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "status": "success",
  "user_id": "12345",
  "conversation_id": "67890",
  "response": "Bonjour ! Bien sûr, je suis là pour vous aider.",
  "timestamp": "2024-11-14T14:32:00Z"
}
```

---

Vous pouvez copier et coller cette description dans un fichier `.md` pour documenter le contrat d'interface de votre API.