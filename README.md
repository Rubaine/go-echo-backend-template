# Go Echo Backend Template

Template minimal de backend en Go utilisant le framework [Echo](https://echo.labstack.com/).
Il est pensé pour être cloné via **degit** puis initialisé rapidement.

## Installation rapide

```bash
npx degit username/go-echo-template mon-projet
cd mon-projet
./init.sh
```

Si vous clonez le dépôt avec `git`, supprimez le dossier `.git` avant d'exécuter `./init.sh`.

## Structure du projet

- `config/` – configuration et accès base de données
- `handlers/` – routes et contrôleurs
- `models/` – structures de données
- `email/` – service d'envoi d'e-mails
- `public/` – fichiers statiques et templates
- `docs/` – documentation Swagger générée
- `Dockerfile`, `docker-compose.yml`

## Utilisation du script d'init

Le script `init.sh` vous demande le chemin de module Go (par ex. `github.com/username/mon-projet`).
Il met à jour `go.mod` et tous les imports, puis lance `go mod tidy`.

```bash
chmod +x init.sh
./init.sh
```

## Technos utilisées

- [Echo](https://echo.labstack.com/)
- [pgx](https://github.com/jackc/pgx) pour PostgreSQL
- [Docker](https://www.docker.com/) & docker-compose
- [Swag](https://github.com/swaggo/swag) pour Swagger

## Licence

Distribué sous licence MIT.
