# Installation du Bot

## Dépendance

Il vous faudra avoir `go` en version 1.20.

## Environnement

L'environnement est à mettre dans `rsa-tcp-support-2023/bot/.env` en suivant le format de `.env.example`.

- `HOST` : IP vers le Hub (127.0.0.1:8000)
- `AUTH` : Mot de passe à envoyer au Hub
- `MOD` & `MOD_R` : Le robot ne répondra que au client dont l'identifiant vérifie `ID%MOD==MOD_R`
- `NAME`: Le nom du robot (que les clients verront)

## Lancement

```bash
git clone https://github.com/yyewolf/rsa-tcp-support-2023/
cd rsa-tcp-support-2023/bot
go run cmd/bot/bot.go
```

Le bot est maintenant lancé, son processus va répondre à tout les messages lui correspondant.