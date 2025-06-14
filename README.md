# Go Backend Template

Template de backend Go moderne avec authentification, API REST, base de données PostgreSQL et Docker. Parfait pour démarrer rapidement un nouveau projet backend.

## 🚀 Démarrage rapide

### Avec degit (recommandé)

```bash
# Installer degit si ce n'est pas déjà fait
npm install -g degit

# Cloner le template
degit your-username/backend-template my-new-project
cd my-new-project

# Initialiser le projet
chmod +x init-template.sh
./init-template.sh my-project-name my-module-name
```

### Avec Git

```bash
git clone https://github.com/your-username/backend-template.git my-new-project
cd my-new-project
rm -rf .git

# Windows
init-template.bat my-project-name my-module-name

# Linux/Mac
chmod +x init-template.sh
./init-template.sh my-project-name my-module-name
```

## 📋 Fonctionnalités incluses

- **🔐 Authentification complète** : inscription, connexion, récupération de mot de passe, gestion des tokens JWT
- **📧 Système d'emails** : envoi d'emails via SMTP avec templates HTML
- **🗄️ Base de données** : PostgreSQL avec migration automatique
- **📚 Documentation** : Swagger/OpenAPI intégré
- **🐳 Docker** : Containerisation complète avec docker-compose
- **🔧 Middleware** : CORS, logging, error handling
- **✅ Structure modulaire** : handlers, models, config séparés
- **🛡️ Sécurité** : hashage des mots de passe, validation des entrées

## 🏗️ Architecture

```
├── config/           # Configuration et variables d'environnement
├── docs/            # Documentation Swagger générée automatiquement
├── email/           # Service d'envoi d'emails avec templates
├── handlers/        # Routes et controllers de l'API
│   └── authHandler/ # Endpoints d'authentification
├── models/          # Structures de données et accès base de données
│   ├── postgresql/  # Configuration PostgreSQL
│   └── user/        # Modèles utilisateur
├── public/          # Assets statiques et templates emails
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── main.go         # Point d'entrée de l'application
```

## ⚙️ Configuration

1. **Copiez le fichier de configuration** :

   ```bash
   cp .env.example .env
   ```

2. **Configurez vos variables dans `.env`** :
   - Base de données PostgreSQL
   - Configuration SMTP pour les emails
   - Secrets JWT
   - Port du serveur

## 🐳 Démarrage avec Docker

```bash
# Démarrer tous les services (API + PostgreSQL)
docker-compose up -d

# Voir les logs
docker-compose logs -f

# Arrêter les services
docker-compose down
```

## 🔧 Démarrage en développement

```bash
# Installer les dépendances
go mod tidy

# Démarrer le serveur
go run main.go

# Ou avec hot reload (installation de air requise)
air
```

## 📖 API Documentation

Une fois le serveur démarré, la documentation Swagger est disponible à :

- **http://localhost:8080/swagger/index.html**

## 🛠️ Endpoints principaux

### Authentification

- `POST /auth/signup` - Inscription
- `POST /auth/login` - Connexion
- `POST /auth/logout` - Déconnexion
- `GET /auth/me` - Profil utilisateur
- `POST /auth/recover` - Demande de récupération de mot de passe
- `POST /auth/reset-password` - Réinitialisation du mot de passe
- `DELETE /auth/signout` - Suppression du compte

### Système

- `GET /health` - Vérification de santé de l'API

## 🔒 Sécurité

- Hashage des mots de passe avec bcrypt
- Authentification JWT avec refresh tokens
- Validation des entrées utilisateur
- Protection CORS configurée
- Middleware de logging des requêtes

## 📦 Dépendances principales

- **Echo v4** - Framework web rapide et minimaliste
- **pgx/v4** - Driver PostgreSQL performant
- **Swaggo** - Génération automatique de documentation Swagger
- **UUID** - Génération d'identifiants uniques
- **Godotenv** - Gestion des variables d'environnement

## 🤝 Contribution

1. Forkez le projet
2. Créez une branche pour votre fonctionnalité
3. Commitez vos changements
4. Poussez vers la branche
5. Ouvrez une Pull Request

## 📄 Licence

Ce template est distribué sous licence MIT. Voir `LICENSE` pour plus d'informations.

- `FRONT_URL` – URL de l'interface web
- `LOG_LEVEL` (par défaut `info`)
- `LISTEN_PORT` (port HTTP, `5000` en développement)
- `MAX_BODY_SIZE` (taille maximale des requêtes, ex. `100M`)
- `MISTRAL_URL`, `MISTRAL_API_KEY`, `MISTRAL_MODEL`

Un fichier `.env.example` fournit un exemple de configuration.

## Lancer l'application

```bash
# Compilation
$ go build -o robert
# Démarrage
$ ./robert
```

En développement il est possible d'utiliser `docker-compose` qui démarre également PostgreSQL :

```bash
$ docker compose up --build
```

L'API écoute par défaut sur le port configuré (`LISTEN_PORT`).

## Principales routes HTTP

| Méthode               | Chemin                 | Description                                |
| --------------------- | ---------------------- | ------------------------------------------ |
| `POST`                | `/auth/signup`         | Création de compte                         |
| `POST`                | `/auth/login`          | Connexion par mot de passe ou token        |
| `POST`                | `/auth/logout`         | Révocation du token                        |
| `POST`                | `/auth/recover`        | Envoi d'un mail de réinitialisation        |
| `POST`                | `/auth/reset_password` | Changement de mot de passe via token       |
| `GET`                 | `/auth/me`             | Informations de l'utilisateur              |
| `POST`                | `/auth/me`             | Mise à jour des informations               |
| `POST`                | `/chat/query`          | Question au LLM avec recherche de services |
| `POST`                | `/chat/page/analyze`   | Analyse d'une page HTML                    |
| `POST`                | `/chat/page/resume`    | Résumé d'une page HTML                     |
| `POST`                | `/chat/mail`           | Analyse d'un e‑mail HTML                   |
| `GET`/`POST`/`DELETE` | `/chat/sessions` ...   | Gestion des sessions de chat               |

## Base de données

Au démarrage, le package `config.InitPgSQL` crée les tables nécessaires :

- `account` : utilisateurs
- `chat_session` et `chat_messages` : historique des discussions

## Prétraitement des services UPHF

Les scripts du dossier `preprocessing/services` permettent de parser le catalogue des services UPHF et d'enrichir chaque service avec des mots‑clés. Le fichier généré `services_augmented.json` est embarqué dans l'image Docker et chargé au lancement pour la recherche sémantique.

```bash
$ cd preprocessing/services
$ python3 parse_uphf_services.py -i services.html --json services.json --csv services.csv
$ python3 augment_services.py -i services.json -o services_augmented.json
```

## Déploiement Docker

Une image minimaliste est construite à partir de `Dockerfile`. Elle embarque l'exécutable compilé, les prompts et les données de services.

```bash
$ docker build -t robert-backend .
$ docker run -p 5000:80 --env-file .env robert-backend
```

## Licence

Ce projet est fourni sous licence MIT.
