# Mori

## Installation Environnement

### Poetry: 

#### Installation:
--> run `pipx install poetry`
--> run `poetry --version` pour voir si Poetry est installée
--> run `poetry install` à la racine du projet
--> run `poetry show` pour voir si les dépendances ont été installées

#### Installation avec Curl:
--> faire dans le terminal `curl -sSL https://install.python-poetry.org | python3 -`
--> Ajoutez cette ligne à votre fichier ~/.bashrc, ~/.zshrc:
    --> `export PATH="$HOME/.local/bin:$PATH"`
--> puis run `source ~/.bashrc` ou `source ~/.zshrc`
--> run `poetry --version` pour voir si poetry est bien installé
--> run `poetry install` à la racine du projet
--> run `poetry show` pour voir si les dépendances ont été installées




### Llama3.1:8b :
--> Aller dans le terminal est run `ollama pull llama3.1:8b` pour download le model en local

## Lancer le server python FastAPI pour le llm: 

--> taper `cd backend/logic_llm`
--> ensuite lancer la commande suivante `uvicorn server:app --host 127.0.0.1 --port 8000`
--> le server FastAPI devrait se lancer.
