# Installation du Hub

## Dépendances

Il vous faudra avoir `go` en version 1.20.

## Environnement

L'environnement est à mettre dans `rsa-tcp-support-2023/hub/.env` en suivant le format de `.env.example`.

- `AUTH_SECRET_0` : mot de passe pour les agents de niveau 1
- `AUTH_SECRET_1` : mot de passe pour les agents de niveau 2
- `AUTH_SECRET_2` : mot de passe pour les agents de niveau 3
- `DEBUG` : true ou false 

## Lancement

```bash
git clone https://github.com/yyewolf/rsa-tcp-support-2023/
cd rsa-tcp-support-2023/hub
go run cmd/hub/hub.go
```

Il faudra laisser tourner le hub pour le bon fonctionnement du reste, celui-ci est sur `:8000`