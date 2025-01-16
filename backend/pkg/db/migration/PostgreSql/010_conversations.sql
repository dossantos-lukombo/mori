-- +migrate Up
CREATE TABLE IF NOT EXISTS conversations (
    id SERIAL PRIMARY KEY, -- Identifiant unique de la conversation
    user_id VARCHAR(100) NOT NULL, -- Identifiant de l'utilisateur, clé étrangère vers la table users
    conversation_id VARCHAR(255) NOT NULL, -- Identifiant unique de la conversation
    user_request TEXT NOT NULL, -- Requête de l'utilisateur
    llm_response TEXT NOT NULL, -- Réponse du LLM
    new_conversation BOOLEAN DEFAULT FALSE NOT NULL, -- Indique si c'est une nouvelle conversation
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Date de création de l'enregistrement
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Dernière mise à jour
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- +migrate Down

DROP TABLE IF EXISTS conversations;
